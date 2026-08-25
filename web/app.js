const metrics = document.querySelector("#metrics");
const reminders = document.querySelector("#reminders");
const records = document.querySelector("#records");

function renderMetrics(snapshot) {
  metrics.replaceChildren(...[
    ["活跃客户", snapshot.TotalClients],
    ["回访记录", snapshot.TotalRecords],
    ["平均满意度", snapshot.AverageScore.toFixed(2)]
  ].map(([label, value]) => {
    const item = document.createElement("article");
    item.innerHTML = `<strong>${value}</strong><span>${label}</span>`;
    return item;
  }));
}

function renderReminders(items) {
  reminders.innerHTML = items.length ? items.map(item => `<p><b>${item.ClientName}</b> ${item.Reason} · ${item.DueDate.slice(0, 10)}</p>`).join("") : "<p>暂无提醒</p>";
}

function renderRecords(items) {
  records.innerHTML = items.map(item => `<article><b>${item.ClientName}</b><span>${item.ServiceTypeName} · ${item.CaregiverName}</span><span>满意度 ${item.Satisfaction.Score}/5 · ${item.Status}</span></article>`).join("");
}

async function load() {
  const [dashboard, recordList] = await Promise.all([fetch("/api/dashboard"), fetch("/api/records")]);
  const snapshot = await dashboard.json();
  const items = await recordList.json();
  renderMetrics(snapshot);
  renderReminders(snapshot.Reminders);
  renderRecords(items);
}

load().catch(error => { reminders.textContent = `加载失败: ${error.message}`; });
