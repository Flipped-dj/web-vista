<template>
    <el-form ref="ruleFormRef" :model="formData">
        <el-row :gutter="24">
            <template v-for="item in formItemAttrs" :key="item.prop" >
                <el-col v-bind="item.col">
                    <el-form-item :label="item.label" :prop="item.prop">
                        <component v-model="formData[item.prop]" :is="isComp(item.comp)" :placeholder="item.placeholder">
                            <template v-if="item.comp === 'select'">
                                <el-option label="全部" value=""></el-option>
                                <el-option 
                                    v-for="opt in item.options" 
                                    :key="opt.value" 
                                    :label="opt.label" 
                                    :value="opt.value"/> 
                            </template>
                        </component>
                    </el-form-item>
                </el-col>
            </template>
        </el-row>
        <el-row>
            <el-button type="primary" @click="handleSearch">查询</el-button>
            <el-button @click="handleReset(ruleFormRef)">重置</el-button>
        </el-row>
    </el-form>
</template>

<script setup>
import { ref,reactive ,computed} from 'vue'

const props = defineProps({
    formItem: {
        type: Array,
        default: () => []
    }
})
const emit = defineEmits(['search'])

const formItemAttrs=computed(()=>{
    const {formItem}=props
    formItem.forEach(item=>{
        item.col={xs:24,sm:12,md:8,lg:6,xl:4}
    })
    return formItem
})

//定义表单响应式数据
const ruleFormRef = ref()
const formData = reactive({})

//动态绑定组件
const isComp=(comp)=>{
    return {
        input:'elInput',
        select:'elSelect',
    }[comp]
}

//表单的查询和重置事件
const handleSearch=()=>{
    emit('search',formData)
}
const handleReset=(formEl)=>{
    //先重置表单，然后再触发查询
    if(!formEl) return
    formEl.resetFields()
    emit('search',formData)
}
</script>
