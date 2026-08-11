<template>
  <el-dialog
    :close-on-click-modal="false"
    width="900px"
    top="5vh"
    v-model="dialogVisible"
    title="考生信息"
    :append-to-body="true"
  >
    <div class="dialog-content">
      <div class="photo">
        <h3>考生照片：</h3>
        <div class="photo-image">
          <el-image :src="dialogProps.row.renZhengZP" fit="contain" />
        </div>
      </div>
      <h3>考生信息：</h3>
      <el-form ref="editFormRef" class="form" label-suffix=" :" :model="dialogProps.row" label-width="130px" disabled>
        <el-row>
          <el-col :span="12">
            <el-form-item label="姓名">
              <el-input v-model="dialogProps.row.xingMing" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="文理科">
              <el-input v-model="dialogProps.row.wenLiKe" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="证件类型">
              <el-input v-model="dialogProps.row.zhengJianLXDesc" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="首选科目">
              <el-input v-model="dialogProps.row.singleSubjectName" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="证件号码">
              <el-input v-model="dialogProps.row.shenFenZH" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="再选科目">
              <el-input v-model="dialogProps.row.subjectReChoose" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="准考证号">
              <el-input v-model="dialogProps.row.zhunKaoZH" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="选考科目">
              <el-input v-model="dialogProps.row.subjectChoosen" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="高考/学籍省份">
              <el-input v-model="dialogProps.row.shengFenMC" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="专业课学习学校">
              <el-input v-model="dialogProps.row.suoZaiHS" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="生日">
              <el-input v-model="dialogProps.row.chuShengRQ" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="文化课学习学校">
              <el-input v-model="dialogProps.row.suoZaiXX" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="学历">
              <el-input v-model="dialogProps.row.xueLi" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="通讯地址">
              <el-input v-model="dialogProps.row.tongXinDZ" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="性别">
              <el-input v-model="dialogProps.row.xingBie" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="民族">
              <el-input v-model="dialogProps.row.minZu" />
            </el-form-item>
          </el-col>
          <!-- <el-col :span="12">
          <el-form-item label="收件人">
            <el-input v-model="dialogProps.row.phone" />
          </el-form-item>
        </el-col> -->
          <el-col :span="12">
            <el-form-item label="政治面貌">
              <el-input v-model="dialogProps.row.zhengZhiMM" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="手机">
              <el-input v-model="dialogProps.row.shouJi" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="考生类型">
              <el-input v-model="dialogProps.row.stuTypeStr" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="应往届">
              <el-input v-model="dialogProps.row.yingWangJie" />
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>
      <h3>报考专业：</h3>
      <el-table :data="tableData" style="width: 100%" border>
        <el-table-column prop="xueXiaoMC" label="院校" align="center" width="160" />
        <el-table-column prop="kaoDianMC" label="考点" align="center" width="150" />
        <el-table-column prop="zhuanYeMC" label="专业" align="center" />
        <!-- <el-table-column prop="address" label="志愿" align="center" width="180" /> -->
        <el-table-column prop="kaoShiRQSM" label="考试时间" align="center" width="100" />
        <el-table-column prop="zhunKaoZH" label="准考证号" align="center" width="130" />
        <el-table-column prop="KaoChangMC" label="考场" align="center" width="180" />
      </el-table>
    </div>
  </el-dialog>
</template>
<script setup lang="ts" name="ViewDialog">
import { ref } from "vue";

interface DialogProps {
  row: any;
}

const dialogVisible = ref(false);
const dialogProps = ref<DialogProps>({
  row: {}
});

const tableData = ref<any[]>([]);

// 接收父组件传过来的参数
const acceptParams = async (params: DialogProps) => {
  dialogProps.value = params;
  dialogVisible.value = true;
  tableData.value = [params.row];
};

defineExpose({
  acceptParams
});
</script>
<style lang="scss" scoped>
.dialog-content {
  height: 500px;
  padding-right: 10px;
  overflow: scroll;
}
.photo {
  display: flex;
  align-items: center;
  margin-bottom: 20px;
  .photo-image {
    width: 150px;
    height: 150px;
    margin-left: 20px;
    overflow: hidden;
    border: 1px solid #dcdfe6;
    border-radius: 4px;
    .el-image {
      width: 100%;
      height: 100%;
    }
  }
}
</style>
