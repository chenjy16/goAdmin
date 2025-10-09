// 调试模型数据流的脚本
import axios from 'axios';

async function debugModelFlow() {
  try {
    console.log('=== 调试模型数据流 ===');
    
    // 1. 测试API端点
    console.log('\n1. 测试API端点...');
    const response = await axios.get('http://localhost:8080/api/v1/ai/googleai/models');
    console.log('API响应状态:', response.status);
    console.log('API响应数据:', JSON.stringify(response.data, null, 2));
    
    // 2. 模拟前端数据转换
    console.log('\n2. 模拟前端数据转换...');
    const modelsObject = response.data.data.models;
    console.log('模型对象:', JSON.stringify(modelsObject, null, 2));
    
    const modelsArray = Object.values(modelsObject);
    console.log('转换后的模型数组:', JSON.stringify(modelsArray, null, 2));
    console.log('模型数组长度:', modelsArray.length);
    
    // 3. 检查每个模型的字段
    console.log('\n3. 检查模型字段...');
    modelsArray.forEach((model, index) => {
      console.log(`模型 ${index + 1}:`, {
        name: model.name,
        display_name: model.display_name,
        enabled: model.enabled,
        max_tokens: model.max_tokens
      });
    });
    
  } catch (error) {
    console.error('错误:', error.message);
    if (error.response) {
      console.error('响应状态:', error.response.status);
      console.error('响应数据:', error.response.data);
    }
  }
}

debugModelFlow();