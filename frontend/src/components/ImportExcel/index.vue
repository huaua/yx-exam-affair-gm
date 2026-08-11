<template>
  <el-dialog
    v-model="dialogVisible"
    :title="`${parameter.title}`"
    :close-on-click-modal="false"
    :destroy-on-close="true"
    width="580px"
    :before-close="handleClose"
    draggable
    center
  >
    <div v-if="parameter.selectedRoomTaskObj" style="padding: 0 25px">
      <p>考场征集任务：{{ parameter.selectedRoomTaskObj?.label || "--" }}</p>
      <p>日期：{{ parameter.selectedRoomTaskObj?.date || "--" }}</p>
    </div>
    <el-form class="drawer-multiColumn-form" label-width="100px">
      <el-form-item label="模板下载 :">
        <el-button type="primary" :icon="Download">
          <!-- <a :href="parameter.templateUrl" target="_blank" download>点击下载</a> -->
          <a :href="handleDownload(parameter.templateUrl)" target="_blank" download>点击下载</a>
        </el-button>
      </el-form-item>
      <el-form-item label="文件上传 :">
        <el-upload
          class="upload"
          ref="uploadRef"
          action="#"
          :drag="true"
          :limit="excelLimit"
          :multiple="true"
          :show-file-list="true"
          :http-request="uploadExcel"
          :on-remove="handleRemove"
          :before-upload="beforeExcelUpload"
          :on-exceed="handleExceed"
          :on-success="excelUploadSuccess"
          :accept="parameter.fileType!.join(',')"
        >
          <slot name="empty">
            <el-icon class="el-icon--upload"><upload-filled /></el-icon>
            <div class="el-upload__text">将文件拖到此处，或<em>点击上传</em></div>
          </slot>
          <template #tip>
            <slot name="tip">
              <div class="el-upload__tip">请上传 .xls , .xlsx 标准格式文件，文件最大为 {{ parameter.fileSize }}M</div>
            </slot>
          </template>
        </el-upload>
      </el-form-item>
    </el-form>
    <template #footer>
      <el-button @click="cancellation">取消</el-button>
      <el-button type="primary" @click="save">上传</el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts" name="ImportExcel">
import { TaskEnumItem } from "@/hooks/useCustomEnum";
import { ref } from "vue";
import { Download } from "@element-plus/icons-vue";
import { upLoad } from "@/api/modules/common";
import { ElMessage, ElNotification, UploadRequestOptions, UploadRawFile, UploadInstance, ElMessageBox } from "element-plus";

export interface ExcelParameterProps {
  title: string; // 标题
  fileSize?: number; // 上传文件的大小
  fileType?: File.ExcelMimeType[]; // 上传文件的类型
  selectedRoomTaskObj?: TaskEnumItem;
  tempApi?: (params: any) => Promise<any>; // 下载模板的Api (暂时未用到)
  importApi?: (params: any) => Promise<any>; // 批量导入的Api （暂时未用到）
  templateUrl?: string;
  saveUploadApi?: (params: any) => Promise<any>; // 保存上传的Api
  getTableList?: () => void; // 获取表格数据的Api
  extraData?: Record<string, any>; // 业务接口附加参数
}

interface MoreFileData {
  originalName?: string;
  url?: string;
}

const uploadRef = ref<UploadInstance>();

// 最大文件上传数
const excelLimit = ref(1);
// dialog状态
const dialogVisible = ref(false);
// 接收选择要上传的文件
const uploadedFiles = ref<any>([]);
// 多文件上传需要使用的参数，需要增加 上传 按钮
let moreFile = ref();
// 获取后端上传文件接口返回的数据
let moreFileData = ref<MoreFileData>({});
// 保存数据接口参数
let filesObj = ref<any>({});

// 父组件传过来的参数
const parameter = ref<ExcelParameterProps>({
  title: "",
  fileSize: 5,
  fileType: ["application/vnd.ms-excel", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"],
  importApi: upLoad
});

// 接收父组件参数
const acceptParams = (params: ExcelParameterProps) => {
  parameter.value = { ...parameter.value, ...params };
  console.log("接收父组件导入的参数 :>> ", parameter.value);
  dialogVisible.value = true;
};

const handleDownload = (url: any) => {
  console.log(`${url}`);
  // console.log(" window.location.host :>> ", window.location.host);
  // return "http://" + window.location.host + url;
  return import.meta.env.VITE_PUBLIC_PATH + url;
};

// 自动文件上传
const uploadExcel = async (param: UploadRequestOptions) => {
  console.log("param :>> ", param);
  uploadedFiles.value.push(param);
  let excelFormData = new FormData();
  //手动添加表单，将文件追加到表单里
  for (let i = 0; i < uploadedFiles.value.length; i++) {
    excelFormData.append("file", uploadedFiles.value[i].file);
  }
  console.log(" 单文件上传 ==> 直接自动请求上传接口 ");
  const res = await parameter.value.importApi!(excelFormData);
  moreFileData.value = res;
};

// 手动删除文件
const handleRemove = () => {
  uploadedFiles.value = [];
  moreFileData.value = {};
};

// 数据重置
const resetData = () => {
  uploadedFiles.value = []; // 上传选择的文件
  excelLimit.value = 1; // 上传文件最大数量
  moreFile.value = ""; // 请求上传接口的参数
  moreFileData.value = {}; // 获取上传接口返回的数据
  filesObj.value = {}; // 保存的参数
  uploadRef.value!.clearFiles(); // 清空所有上传的文件
};

// 上传
const save = async () => {
  try {
    console.log("success save!");
    if (!uploadedFiles.value.length) {
      ElMessage({
        message: `请上传文件！`,
        type: "warning"
      });
      return;
    }
    const fileNameList = moreFileData.value.url?.split(",");
    console.log("fileNameList :>> ", fileNameList);
    const originalNameList = moreFileData.value.originalName?.split(",");
    console.log("originalNameList :>> ", originalNameList);
    originalNameList?.forEach((item, index) => {
      filesObj.value[item] = fileNameList?.length ? fileNameList[fileNameList.length - 1 - index] : "";
    });
    console.log("filesObj :>> ", filesObj.value);

    const dataObj = {
      files: filesObj.value,
      roomTaskId: parameter.value.selectedRoomTaskObj?.value,
      ...(parameter.value.extraData || {})
    };
    await parameter.value.saveUploadApi!(dataObj);
    ElMessage({
      message: `上传成功！`,
      type: "success"
    });
    parameter.value.getTableList && parameter.value.getTableList();
    resetData(); // 重置数据
    dialogVisible.value = false; // 控制弹窗展开收起
  } catch (error) {
    console.log("error save!");
  }
};

// 取消
const cancellation = () => {
  resetData();
  dialogVisible.value = false; // 控制弹窗展开收起
};

// 点击右上角x关闭dialog弹框
const handleClose = (done: () => void) => {
  ElMessageBox.confirm("确定关闭吗?")
    .then(() => {
      resetData();
      console.log("是否执行done :>> ");
      done();
    })
    .catch(() => {
      // catch error
    });
};

/**
 * @description 文件上传之前判断
 * @param file 上传的文件
 * */
const beforeExcelUpload = (file: UploadRawFile) => {
  const isExcel = parameter.value.fileType!.includes(file.type as File.ExcelMimeType);
  const fileSize = file.size / 1024 / 1024 < parameter.value.fileSize!;
  if (!isExcel) ElMessage.warning("上传文件只能是 xls / xlsx 格式！");
  if (!fileSize) ElMessage.warning(`上传文件大小不能超过 ${parameter.value.fileSize}MB！`);
  return isExcel && fileSize;
};

// 文件数超出提示
const handleExceed = () => {
  ElMessage.warning("最多只能上传一个文件！");
};

// 上传成功提示
const excelUploadSuccess = () => {
  ElNotification({
    title: "温馨提示",
    message: `文件导入成功，请点击上传！`,
    type: "success"
  });
};

defineExpose({
  acceptParams
});
</script>
<style lang="scss" scoped>
@import "./index";

/* 去除a标签下划线 */

/* 去除默认的颜色和点击后变化的颜色 */
a {
  color: #ffffff;
  text-decoration: none;
}

/* 去除未被访问的a标签的下划线 */
a:link {
  text-decoration: none;
}

/* 去除已经被访问过的a标签的下划线 */
a:visited {
  text-decoration: none;
}

/* 去除鼠标悬停时的a标签的下划线 */
a:hover {
  text-decoration: none;
}

/* 去除正在点击的a标签的下划线（鼠标按下，尚未松开） */
a:active {
  text-decoration: none;
}

/* 去除获得焦点的a标签的下划线（被鼠标点击过） */
a:focus {
  text-decoration: none;
}
</style>
