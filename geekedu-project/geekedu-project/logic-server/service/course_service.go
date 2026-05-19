package service

import (
	"context"
	"strings"

	"geekedu/common/config"
	"geekedu/common/oss"
	"geekedu/logic-server/database"
	"geekedu/logic-server/model"
	"geekedu/proto/pb"

	"gorm.io/gorm"
)

type CourseService struct {
	pb.UnimplementedCourseServiceServer
	ossConfig *config.OSSConfig
}

func NewCourseService(ossConfig *config.OSSConfig) *CourseService {
	return &CourseService{
		ossConfig: ossConfig,
	}
}

// GetCourseList 获取课程列表
func (s *CourseService) GetCourseList(ctx context.Context, req *pb.GetCourseListRequest) (*pb.GetCourseListResponse, error) {
	var courses []model.Course
	var total int64

	query := database.DB.Model(&model.Course{}).Where("status = 1")

	// 按分类筛选
	if req.CategoryId > 0 {
		query = query.Where("category_id = ?", req.CategoryId)
	}

	// 查询总数
	if err := query.Count(&total).Error; err != nil {
		return &pb.GetCourseListResponse{
			Code:    500,
			Message: "query failed",
		}, err
	}

	// 分页查询
	offset := int((req.PageNum - 1) * req.PageSize)
	if err := query.Offset(offset).Limit(int(req.PageSize)).Find(&courses).Error; err != nil {
		return &pb.GetCourseListResponse{
			Code:    500,
			Message: "query failed",
		}, err
	}

	// 转换为proto格式
	list := make([]*pb.CourseInfo, len(courses))
	for i, course := range courses {
		list[i] = &pb.CourseInfo{
			Id:            course.ID,
			Title:         course.Title,
			Cover:         course.Cover,
			Description:   course.Description,
			CategoryId:    course.CategoryID,
			TeacherId:     course.TeacherID,
			Price:         course.Price,
			OriginalPrice: course.OriginalPrice,
			StudentCount:  int32(course.StudentCount),
			LessonCount:   int32(course.LessonCount),
			Duration:      int32(course.Duration),
			Level:         int32(course.Level),
		}
	}

	return &pb.GetCourseListResponse{
		Code:    0,
		Message: "success",
		List:    list,
		Total:   total,
	}, nil
}

// GetCourseDetail 获取课程详情
func (s *CourseService) GetCourseDetail(ctx context.Context, req *pb.GetCourseDetailRequest) (*pb.GetCourseDetailResponse, error) {
	var course model.Course
	err := database.DB.Where("id = ?", req.CourseId).First(&course).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &pb.GetCourseDetailResponse{
				Code:    404,
				Message: "course not found",
			}, nil
		}
		return &pb.GetCourseDetailResponse{
			Code:    500,
			Message: "query failed",
		}, err
	}

	// 查询章节和视频
	var chapters []model.Chapter
	database.DB.Where("course_id = ?", course.ID).Order("sort").Find(&chapters)

	chapterInfos := make([]*pb.ChapterInfo, len(chapters))
	for i, chapter := range chapters {
		var videos []model.Video
		database.DB.Where("chapter_id = ?", chapter.ID).Order("sort").Find(&videos)

		videoInfos := make([]*pb.VideoInfo, len(videos))
		for j, video := range videos {
			videoInfos[j] = &pb.VideoInfo{
				Id:       video.ID,
				Title:    video.Title,
				Duration: int32(video.Duration),
				VideoUrl: video.VideoURL,
				Cover:    video.Cover,
				IsFree:   int32(video.IsFree),
			}
		}

		chapterInfos[i] = &pb.ChapterInfo{
			Id:     chapter.ID,
			Title:  chapter.Title,
			Sort:   int32(chapter.Sort),
			Videos: videoInfos,
		}
	}

	// 查询讲师信息
	var teacher model.User
	database.DB.Where("id = ?", course.TeacherID).First(&teacher)

	return &pb.GetCourseDetailResponse{
		Code:    0,
		Message: "success",
		CourseDetail: &pb.CourseDetail{
			CourseInfo: &pb.CourseInfo{
				Id:            course.ID,
				Title:         course.Title,
				Cover:         course.Cover,
				Description:   course.Description,
				CategoryId:    course.CategoryID,
				TeacherId:     course.TeacherID,
				Price:         course.Price,
				OriginalPrice: course.OriginalPrice,
				StudentCount:  int32(course.StudentCount),
				LessonCount:   int32(course.LessonCount),
				Duration:      int32(course.Duration),
				Level:         int32(course.Level),
			},
			Chapters: chapterInfos,
			Teacher: &pb.TeacherInfo{
				Id:       teacher.ID,
				Username: teacher.Username,
				Nickname: teacher.Nickname,
				Avatar:   "",
			},
		},
	}, nil
}

// GetCategoryTree 获取分类树
func (s *CourseService) GetCategoryTree(ctx context.Context, req *pb.GetCategoryTreeRequest) (*pb.GetCategoryTreeResponse, error) {
	var categories []model.Category
	if err := database.DB.Order("sort").Find(&categories).Error; err != nil {
		return &pb.GetCategoryTreeResponse{
			Code:    500,
			Message: "query failed",
		}, err
	}

	// 获取每个分类的课程数量
	type CategoryCount struct {
		CategoryID int64
		Count      int32
	}
	var counts []CategoryCount
	database.DB.Model(&model.Course{}).
		Select("category_id, count(*) as count").
		Where("status = 1").
		Group("category_id").
		Scan(&counts)

	countMap := make(map[int64]int32)
	for _, c := range counts {
		countMap[c.CategoryID] = c.Count
	}

	// 分类映射
	categoryMap := make(map[int64]*pb.CategoryInfo)
	var rootCategories []*pb.CategoryInfo

	// 第一次遍历：创建所有节点
	for _, cat := range categories {
		categoryMap[cat.ID] = &pb.CategoryInfo{
			Id:          cat.ID,
			Name:        cat.Name,
			ParentId:    cat.ParentID,
			Children:    []*pb.CategoryInfo{},
			CourseCount: countMap[cat.ID],
			Sort:        int32(cat.Sort),
		}
	}

	// 第二次遍历：构建树形结构
	for _, cat := range categories {
		node := categoryMap[cat.ID]
		if cat.ParentID == 0 {
			rootCategories = append(rootCategories, node)
		} else {
			if parent, exists := categoryMap[cat.ParentID]; exists {
				parent.Children = append(parent.Children, node)
			}
		}
	}

	return &pb.GetCategoryTreeResponse{
		Code:    0,
		Message: "success",
		List:    rootCategories,
	}, nil
}

// CreateCourse 创建课程
func (s *CourseService) CreateCourse(ctx context.Context, req *pb.CreateCourseRequest) (*pb.CreateCourseResponse, error) {
	course := &model.Course{
		Title:       req.Title,
		Description: req.Description,
		Cover:       req.CoverUrl,
		Price:       req.Price,
		CategoryID:  req.CategoryId,
		TeacherID:   req.TeacherId,
		Status:      0, // 待审核
	}

	if err := database.DB.Create(course).Error; err != nil {
		return &pb.CreateCourseResponse{
			Code:    500,
			Message: "",
		}, err
	}

	return &pb.CreateCourseResponse{
		Code:     0,
		Message:  "course created, pending review",
		CourseId: course.ID,
	}, nil
}

// UpdateCourse 更新课程
func (s *CourseService) UpdateCourse(ctx context.Context, req *pb.UpdateCourseRequest) (*pb.UpdateCourseResponse, error) {
	updates := map[string]interface{}{}
	if req.Title != "" {
		updates["title"] = req.Title
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.CoverUrl != "" {
		updates["cover"] = req.CoverUrl
	}
	if req.Price > 0 {
		updates["price"] = req.Price
	}
	if req.CategoryId > 0 {
		updates["category_id"] = req.CategoryId
	}

	err := database.DB.Model(&model.Course{}).Where("id = ?", req.CourseId).Updates(updates).Error
	if err != nil {
		return &pb.UpdateCourseResponse{
			Code:    500,
			Message: "",
		}, err
	}

	return &pb.UpdateCourseResponse{
		Code:    0,
		Message: "updated",
	}, nil
}

// DeleteCourse 删除课程
func (s *CourseService) DeleteCourse(ctx context.Context, req *pb.DeleteCourseRequest) (*pb.DeleteCourseResponse, error) {
	// 检查课程是否存在且属于该讲师
	var course model.Course
	err := database.DB.Where("id = ?", req.CourseId).First(&course).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &pb.DeleteCourseResponse{Code: 404, Message: "course not found"}, nil
		}
		return &pb.DeleteCourseResponse{Code: 500, Message: "query failed"}, err
	}

	// 检查是否为课程创建者（或管理员）
	if course.TeacherID != req.TeacherId {
		// 检查是否为管理员
		var user model.User
		userErr := database.DB.Where("id = ?", req.TeacherId).First(&user).Error
		if userErr != nil || user.Role != 2 {
			return &pb.DeleteCourseResponse{Code: 403, Message: ""}, nil
		}
	}

	// 删除课程下的视频
	database.DB.Where("course_id = ?", req.CourseId).Delete(&model.Video{})
	// 删除课程下的章节
	database.DB.Where("course_id = ?", req.CourseId).Delete(&model.Chapter{})
	// 删除课程
	err = database.DB.Delete(&course).Error
	if err != nil {
		return &pb.DeleteCourseResponse{Code: 500, Message: ""}, err
	}

	return &pb.DeleteCourseResponse{Code: 0, Message: "deleted"}, nil
}

// UploadVideo 上传视频
func (s *CourseService) UploadVideo(ctx context.Context, req *pb.UploadVideoRequest) (*pb.UploadVideoResponse, error) {
	// 验证chapterId - 必须指定有效的章节
	chapterId := req.ChapterId
	if chapterId == 0 {
		return &pb.UploadVideoResponse{
			Code:    400,
			Message: "please select a chapter",
		}, nil
	}

	// 验证章节存在且属于该课程
	var chapter model.Chapter
	err := database.DB.Where("id = ? AND course_id = ?", chapterId, req.CourseId).First(&chapter).Error
	if err != nil {
		return &pb.UploadVideoResponse{
			Code:    404,
			Message: "chapter not found",
		}, nil
	}

	video := &model.Video{
		CourseID:  req.CourseId,
		ChapterID: &chapterId,
		Title:     req.Title,
		VideoURL:  req.VideoUrl,
		Duration:  int(req.Duration),
		IsFree:    int(req.IsFree),
	}

	if err := database.DB.Create(video).Error; err != nil {
		return &pb.UploadVideoResponse{
			Code:    500,
			Message: "",
		}, err
	}

	//
	database.DB.Model(&model.Course{}).Where("id = ?", req.CourseId).
		UpdateColumn("lesson_count", gorm.Expr("lesson_count + ?", 1))

	return &pb.UploadVideoResponse{
		Code:    0,
		Message: "",
		VideoId: video.ID,
	}, nil
}

// GetCourseVideos 获取课程视频列表
func (s *CourseService) GetCourseVideos(ctx context.Context, req *pb.GetCourseVideosRequest) (*pb.GetCourseVideosResponse, error) {
	var videos []model.Video
	err := database.DB.Where("course_id = ?", req.CourseId).Order("sort").Find(&videos).Error
	if err != nil {
		return &pb.GetCourseVideosResponse{
			Code:    500,
			Message: "query failed",
		}, err
	}

	videoInfos := make([]*pb.VideoInfo, len(videos))
	for i, v := range videos {
		videoInfos[i] = &pb.VideoInfo{
			Id:       v.ID,
			Title:    v.Title,
			Duration: int32(v.Duration),
			VideoUrl: v.VideoURL,
			Cover:    v.Cover,
			IsFree:   int32(v.IsFree),
		}
	}

	return &pb.GetCourseVideosResponse{
		Code:    0,
		Message: "success",
		Videos:  videoInfos,
	}, nil
}

// CreateChapter 创建章节
func (s *CourseService) CreateChapter(ctx context.Context, req *pb.CreateChapterRequest) (*pb.CreateChapterResponse, error) {
	// 验证课程是否存在
	var course model.Course
	if err := database.DB.Where("id = ?", req.CourseId).First(&course).Error; err != nil {
		return &pb.CreateChapterResponse{
			Code:    404,
			Message: "course not found",
		}, nil
	}

	chapter := &model.Chapter{
		CourseID: req.CourseId,
		Title:    req.Title,
		Sort:     int(req.Sort),
	}

	if err := database.DB.Create(chapter).Error; err != nil {
		return &pb.CreateChapterResponse{
			Code:    500,
			Message: "create chapter failed",
		}, err
	}

	return &pb.CreateChapterResponse{
		Code:      0,
		Message:   "chapter created",
		ChapterId: chapter.ID,
	}, nil
}

// UpdateCourseStatus 更新课程状态
func (s *CourseService) UpdateCourseStatus(ctx context.Context, req *pb.UpdateCourseStatusRequest) (*pb.UpdateCourseStatusResponse, error) {
	err := database.DB.Model(&model.Course{}).Where("id = ?", req.CourseId).Update("status", req.Status).Error
	if err != nil {
		return &pb.UpdateCourseStatusResponse{
			Code:    500,
			Message: "",
		}, err
	}

	statusNames := []string{"pending", "published", "unpublished"}
	return &pb.UpdateCourseStatusResponse{
		Code:    0,
		Message: "course status updated to: " + statusNames[req.Status],
	}, nil
}

// CreateCategory 创建分类
func (s *CourseService) CreateCategory(ctx context.Context, req *pb.CreateCategoryRequest) (*pb.CreateCategoryResponse, error) {
	category := &model.Category{
		Name:     req.Name,
		ParentID: req.ParentId,
		Sort:     int(req.Sort),
	}

	if err := database.DB.Create(category).Error; err != nil {
		return &pb.CreateCategoryResponse{
			Code:    500,
			Message: "",
		}, err
	}

	return &pb.CreateCategoryResponse{
		Code:       0,
		Message:    "created",
		CategoryId: category.ID,
	}, nil
}

// UpdateCategory 更新分类
func (s *CourseService) UpdateCategory(ctx context.Context, req *pb.UpdateCategoryRequest) (*pb.UpdateCategoryResponse, error) {
	updates := map[string]interface{}{
		"name": req.Name,
		"sort": req.Sort,
	}

	err := database.DB.Model(&model.Category{}).Where("id = ?", req.CategoryId).Updates(updates).Error
	if err != nil {
		return &pb.UpdateCategoryResponse{
			Code:    500,
			Message: "",
		}, err
	}

	return &pb.UpdateCategoryResponse{
		Code:    0,
		Message: "updated",
	}, nil
}

// DeleteCategory 删除分类
func (s *CourseService) DeleteCategory(ctx context.Context, req *pb.DeleteCategoryRequest) (*pb.DeleteCategoryResponse, error) {
	// 检查是否有子分类
	var count int64
	database.DB.Model(&model.Category{}).Where("parent_id = ?", req.CategoryId).Count(&count)
	if count > 0 {
		return &pb.DeleteCategoryResponse{
			Code:    400,
			Message: "cannot delete category with sub-categories",
		}, nil
	}

	// 检查分类是否被课程使用
	database.DB.Model(&model.Course{}).Where("category_id = ?", req.CategoryId).Count(&count)
	if count > 0 {
		return &pb.DeleteCategoryResponse{
			Code:    400,
			Message: "",
		}, nil
	}

	err := database.DB.Delete(&model.Category{}, req.CategoryId).Error
	if err != nil {
		return &pb.DeleteCategoryResponse{
			Code:    500,
			Message: "",
		}, err
	}

	return &pb.DeleteCategoryResponse{
		Code:    0,
		Message: "deleted",
	}, nil
}

// GetVideoPlayURL 获取视频播放URL（带鉴权）
func (s *CourseService) GetVideoPlayURL(ctx context.Context, req *pb.GetVideoPlayURLRequest) (*pb.GetVideoPlayURLResponse, error) {
	// 1. 查询视频信息
	var video model.Video
	err := database.DB.Where("id = ?", req.VideoId).First(&video).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &pb.GetVideoPlayURLResponse{
				Code:    404,
				Message: "video not found",
			}, nil
		}
		return &pb.GetVideoPlayURLResponse{
			Code:    500,
			Message: "query failed",
		}, err
	}

	// 2. 查询课程信息
	var course model.Course
	err = database.DB.Where("id = ?", video.CourseID).First(&course).Error
	if err != nil {
		return &pb.GetVideoPlayURLResponse{
			Code:    500,
			Message: "",
		}, err
	}

	// 生成签名URL的辅助函数
	generateSignedURL := func(videoURL string) (string, error) {
		// 检查是否配置了OSS以及是否为OSS URL
		if s.ossConfig == nil || s.ossConfig.AccessKeyID == "" || !strings.Contains(videoURL, "aliyuncs.com") {
			return videoURL, nil
		}
		ossClient, err := oss.NewOSSClient(
			s.ossConfig.Endpoint,
			s.ossConfig.AccessKeyID,
			s.ossConfig.AccessKeySecret,
			s.ossConfig.BucketName,
		)
		if err != nil {
			return videoURL, err
		}
		// 生成有效3600秒的预签名URL
		return ossClient.GetPresignedURL(videoURL, 3600)
	}

	// 3. 检查是否为免费视频
	// 免费视频直接返回签名URL
	if video.IsFree == 1 {
		signedURL, err := generateSignedURL(video.VideoURL)
		if err != nil {
			return &pb.GetVideoPlayURLResponse{
				Code:    0,
				Message: "success",
				PlayUrl: video.VideoURL, // 返回原始URL
			}, nil
		}
		return &pb.GetVideoPlayURLResponse{
			Code:    0,
			Message: "success",
			PlayUrl: signedURL,
		}, nil
	}

	// 检查用户是否登录
	if req.UserId == 0 {
		return &pb.GetVideoPlayURLResponse{
			Code:    401,
			Message: "",
		}, nil
	}

	// 检查是否为课程创建者
	if course.TeacherID == req.UserId {
		signedURL, _ := generateSignedURL(video.VideoURL)
		return &pb.GetVideoPlayURLResponse{
			Code:    0,
			Message: "success",
			PlayUrl: signedURL,
		}, nil
	}

	// 检查用户是否已购买该课程
	var userCourse model.UserCourse
	err = database.DB.Where("user_id = ? AND course_id = ?", req.UserId, video.CourseID).First(&userCourse).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &pb.GetVideoPlayURLResponse{
				Code:    403,
				Message: "please purchase this course first",
			}, nil
		}
		return &pb.GetVideoPlayURLResponse{
			Code:    500,
			Message: "",
		}, err
	}

	// 4. 生成签名URL并返回
	signedURL, _ := generateSignedURL(video.VideoURL)
	return &pb.GetVideoPlayURLResponse{
		Code:    0,
		Message: "success",
		PlayUrl: signedURL,
	}, nil
}
