// 页面加载完成后获取并显示结果
document.addEventListener("DOMContentLoaded", function () {
    fetch("/data-entry")
      .then(response => response.json())
      .then(data => {
        const resultTable = document.getElementById("resultTable").getElementsByTagName('tbody')[0];
        data.forEach(entry => {
          const row = resultTable.insertRow();
          // 根据后端返回的字段添加单元格
          // entry.field1, entry.field2, entry.field3
        });
      })
      .catch(error => console.error("获取结果时出错：", error));
  });
  