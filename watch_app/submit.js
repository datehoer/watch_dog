function submitTask() {
    const db = document.getElementById("db").value;
    const tables = document.getElementById("tables").value;
    const limit = document.getElementById("limit").value;
    const orderby = document.getElementById("orderby").value;
    const update = document.getElementById("update").checked;
  
    // 构建请求参数对象
    const requestData = {
      db: db,
      tables: tables,
      limit: limit,
      orderby: orderby,
      update: update
    };
  
    // 发起POST请求将任务信息发送给后端
    // 可以使用Fetch API或其他AJAX库
  
    // 清空表单
    document.getElementById("db").value = "";
    document.getElementById("tables").value = "";
    document.getElementById("limit").value = "100";
    document.getElementById("orderby").value = "";
    document.getElementById("update").checked = false;
  }
  