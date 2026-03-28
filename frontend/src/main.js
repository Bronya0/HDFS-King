import './style.css';

// ===== SVG Icons =====
const ICONS = {
  folder: `<svg viewBox="0 0 24 24" fill="currentColor"><path d="M10 4H4c-1.1 0-2 .9-2 2v12c0 1.1.9 2 2 2h16c1.1 0 2-.9 2-2V8c0-1.1-.9-2-2-2h-8l-2-2z"/></svg>`,
  file: `<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/><polyline points="14 2 14 8 20 8"/></svg>`,
  plus: `<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><line x1="12" y1="5" x2="12" y2="19"/><line x1="5" y1="12" x2="19" y2="12"/></svg>`,
  upload: `<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="16 16 12 12 8 16"/><line x1="12" y1="12" x2="12" y2="21"/><path d="M20.39 18.39A5 5 0 0 0 18 9h-1.26A8 8 0 1 0 3 16.3"/></svg>`,
  download: `<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="8 17 12 21 16 17"/><line x1="12" y1="12" x2="12" y2="21"/><path d="M20.88 18.09A5 5 0 0 0 18 9h-1.26A8 8 0 1 0 3 16.3"/></svg>`,
  trash: `<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="3 6 5 6 21 6"/><path d="M19 6l-1 14a2 2 0 0 1-2 2H8a2 2 0 0 1-2-2L5 6"/><path d="M10 11v6"/><path d="M14 11v6"/><path d="M9 6V4a1 1 0 0 1 1-1h4a1 1 0 0 1 1 1v2"/></svg>`,
  folderPlus: `<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M22 19a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h5l2 3h9a2 2 0 0 1 2 2z"/><line x1="12" y1="11" x2="12" y2="17"/><line x1="9" y1="14" x2="15" y2="14"/></svg>`,
  edit: `<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"/><path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"/></svg>`,
  refresh: `<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="23 4 23 10 17 10"/><path d="M20.49 15a9 9 0 1 1-2.12-9.36L23 10"/></svg>`,
  back: `<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><line x1="19" y1="12" x2="5" y2="12"/><polyline points="12 19 5 12 12 5"/></svg>`,
  up: `<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><line x1="12" y1="19" x2="12" y2="5"/><polyline points="5 12 12 5 19 12"/></svg>`,
  connect: `<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M10 13a5 5 0 0 0 7.54.54l3-3a5 5 0 0 0-7.07-7.07l-1.72 1.71"/><path d="M14 11a5 5 0 0 0-7.54-.54l-3 3a5 5 0 0 0 7.07 7.07l1.71-1.71"/></svg>`,
  disconnect: `<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M18.84 12.25l1.72-1.71a5 5 0 0 0-7.07-7.07l-1.72 1.71"/><path d="M5.16 11.75l-1.72 1.71a5 5 0 0 0 7.07 7.07l1.72-1.71"/><line x1="1" y1="1" x2="23" y2="23"/></svg>`,
  server: `<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="2" y="2" width="20" height="8" rx="2" ry="2"/><rect x="2" y="14" width="20" height="8" rx="2" ry="2"/><line x1="6" y1="6" x2="6.01" y2="6"/><line x1="6" y1="18" x2="6.01" y2="18"/></svg>`,
  home: `<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M3 9l9-7 9 7v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z"/></svg>`,
};

// ===== State =====
const state = {
  connections: [],
  activeConnId: null,
  connected: false,
  currentPath: '/',
  files: [],
  selectedFile: null,
  history: [],
  historyIndex: -1,
  loading: false,
};

// ===== Wails Backend Bindings =====
// These will be available as window.go.backend.* after wails generates bindings
function getConnMgr() { return window.go.backend.ConnectionManager; }
function getHdfsSvc() { return window.go.backend.HdfsService; }

// ===== Toast =====
function showToast(message, type = 'info') {
  const container = document.getElementById('toast-container');
  const el = document.createElement('div');
  el.className = `toast ${type}`;
  el.textContent = message;
  container.appendChild(el);
  setTimeout(() => { el.remove(); }, 3000);
}

// ===== Format file size =====
function formatSize(bytes) {
  if (bytes === 0) return '-';
  const units = ['B', 'KB', 'MB', 'GB', 'TB'];
  let i = 0;
  let size = bytes;
  while (size >= 1024 && i < units.length - 1) {
    size /= 1024;
    i++;
  }
  return size.toFixed(i === 0 ? 0 : 1) + ' ' + units[i];
}

// ===== Connection Management =====
async function loadConnections() {
  try {
    state.connections = await getConnMgr().GetConnections();
  } catch (e) {
    state.connections = [];
  }
  renderSidebar();
}

async function addConnection() {
  showModal('新建连接', `
    <div class="form-group">
      <label>连接名称</label>
      <input id="modal-conn-name" placeholder="My HDFS Cluster" />
    </div>
    <div class="form-group">
      <label>NameNode 地址</label>
      <input id="modal-conn-addr" placeholder="namenode:9000" />
    </div>
    <div class="form-group">
      <label>用户名（可选）</label>
      <input id="modal-conn-user" placeholder="hdfs" />
    </div>
  `, async () => {
    const name = document.getElementById('modal-conn-name').value.trim();
    const address = document.getElementById('modal-conn-addr').value.trim();
    const user = document.getElementById('modal-conn-user').value.trim();
    if (!name || !address) {
      showToast('名称和地址不能为空', 'error');
      return;
    }
    const res = await getConnMgr().AddConnection(name, address, user);
    if (res.success) {
      showToast('连接已添加', 'success');
      closeModal();
      await loadConnections();
    } else {
      showToast(res.message, 'error');
    }
  });
}

async function editConnection(id) {
  const conn = state.connections.find(c => c.id === id);
  if (!conn) return;
  showModal('编辑连接', `
    <div class="form-group">
      <label>连接名称</label>
      <input id="modal-conn-name" value="${escapeHtml(conn.name)}" />
    </div>
    <div class="form-group">
      <label>NameNode 地址</label>
      <input id="modal-conn-addr" value="${escapeHtml(conn.address)}" />
    </div>
    <div class="form-group">
      <label>用户名（可选）</label>
      <input id="modal-conn-user" value="${escapeHtml(conn.user || '')}" />
    </div>
  `, async () => {
    const name = document.getElementById('modal-conn-name').value.trim();
    const address = document.getElementById('modal-conn-addr').value.trim();
    const user = document.getElementById('modal-conn-user').value.trim();
    const res = await getConnMgr().UpdateConnection(id, name, address, user);
    if (res.success) {
      showToast('连接已更新', 'success');
      closeModal();
      await loadConnections();
    } else {
      showToast(res.message, 'error');
    }
  });
}

async function deleteConnection(id) {
  if (!confirm('确定要删除此连接吗？')) return;
  const res = await getConnMgr().DeleteConnection(id);
  if (res.success) {
    if (state.activeConnId === id) {
      await disconnectHdfs();
    }
    showToast('连接已删除', 'success');
    await loadConnections();
  } else {
    showToast(res.message, 'error');
  }
}

async function connectToHdfs(id) {
  const conn = state.connections.find(c => c.id === id);
  if (!conn) return;

  // 如果当前已连接，先断开
  if (state.connected) {
    await getHdfsSvc().Disconnect();
  }

  setLoading(true);
  try {
    const res = await getHdfsSvc().Connect(conn.address, conn.user || '');
    if (res.success) {
      state.connected = true;
      state.activeConnId = id;
      state.currentPath = '/';
      state.history = ['/'];
      state.historyIndex = 0;
      showToast(`已连接到 ${conn.name}`, 'success');
      await navigateTo('/');
    } else {
      showToast(res.message, 'error');
    }
  } catch (e) {
    showToast('连接失败: ' + e, 'error');
  }
  setLoading(false);
  renderSidebar();
  renderStatusBar();
}

async function disconnectHdfs() {
  await getHdfsSvc().Disconnect();
  state.connected = false;
  state.activeConnId = null;
  state.files = [];
  state.selectedFile = null;
  state.currentPath = '/';
  state.history = [];
  state.historyIndex = -1;
  renderAll();
  showToast('已断开连接', 'info');
}

// ===== File Operations =====
async function navigateTo(dirPath) {
  if (!state.connected) return;
  setLoading(true);
  try {
    const result = await getHdfsSvc().ListDir(dirPath);
    if (result.error) {
      showToast(result.error, 'error');
    }
    state.currentPath = result.path || dirPath;
    state.files = result.files || [];
    state.selectedFile = null;

    // Update history
    if (state.history[state.historyIndex] !== dirPath) {
      state.history = state.history.slice(0, state.historyIndex + 1);
      state.history.push(dirPath);
      state.historyIndex = state.history.length - 1;
    }
  } catch (e) {
    showToast('加载目录失败: ' + e, 'error');
  }
  setLoading(false);
  renderContent();
  renderStatusBar();
}

function goBack() {
  if (state.historyIndex > 0) {
    state.historyIndex--;
    const path = state.history[state.historyIndex];
    navigateWithoutHistory(path);
  }
}

function goUp() {
  if (state.currentPath === '/') return;
  const parts = state.currentPath.split('/').filter(Boolean);
  parts.pop();
  const parentPath = '/' + parts.join('/');
  navigateTo(parentPath || '/');
}

async function navigateWithoutHistory(dirPath) {
  if (!state.connected) return;
  setLoading(true);
  try {
    const result = await getHdfsSvc().ListDir(dirPath);
    if (result.error) {
      showToast(result.error, 'error');
    }
    state.currentPath = result.path || dirPath;
    state.files = result.files || [];
    state.selectedFile = null;
  } catch (e) {
    showToast('加载目录失败: ' + e, 'error');
  }
  setLoading(false);
  renderContent();
  renderStatusBar();
}

async function refreshDir() {
  await navigateWithoutHistory(state.currentPath);
}

async function createFolder() {
  showModal('新建文件夹', `
    <div class="form-group">
      <label>文件夹名称</label>
      <input id="modal-folder-name" placeholder="new_folder" />
    </div>
  `, async () => {
    const name = document.getElementById('modal-folder-name').value.trim();
    if (!name) {
      showToast('文件夹名称不能为空', 'error');
      return;
    }
    const fullPath = (state.currentPath === '/' ? '/' : state.currentPath + '/') + name;
    const res = await getHdfsSvc().MkDir(fullPath);
    if (res.success) {
      showToast('文件夹已创建', 'success');
      closeModal();
      await refreshDir();
    } else {
      showToast(res.message, 'error');
    }
  });
}

async function uploadFile() {
  if (!state.connected) return;
  setLoading(true);
  try {
    const res = await getHdfsSvc().Upload(state.currentPath);
    if (res.success) {
      showToast(res.message, 'success');
      await refreshDir();
    } else if (res.message !== '未选择文件') {
      showToast(res.message, 'error');
    }
  } catch (e) {
    showToast('上传失败: ' + e, 'error');
  }
  setLoading(false);
}

async function downloadFile() {
  if (!state.selectedFile || state.selectedFile.isDir) return;
  setLoading(true);
  try {
    const res = await getHdfsSvc().Download(state.selectedFile.path);
    if (res.success) {
      showToast(res.message, 'success');
    } else if (res.message !== '未选择保存位置') {
      showToast(res.message, 'error');
    }
  } catch (e) {
    showToast('下载失败: ' + e, 'error');
  }
  setLoading(false);
}

async function deleteSelected() {
  if (!state.selectedFile) return;
  const name = state.selectedFile.name;
  if (!confirm(`确定要删除 "${name}" 吗？此操作不可恢复。`)) return;
  setLoading(true);
  try {
    const res = await getHdfsSvc().Delete(state.selectedFile.path);
    if (res.success) {
      showToast('删除成功', 'success');
      state.selectedFile = null;
      await refreshDir();
    } else {
      showToast(res.message, 'error');
    }
  } catch (e) {
    showToast('删除失败: ' + e, 'error');
  }
  setLoading(false);
}

async function renameSelected() {
  if (!state.selectedFile) return;
  const oldName = state.selectedFile.name;
  const oldPath = state.selectedFile.path;
  showModal('重命名', `
    <div class="form-group">
      <label>新名称</label>
      <input id="modal-rename-name" value="${escapeHtml(oldName)}" />
    </div>
  `, async () => {
    const newName = document.getElementById('modal-rename-name').value.trim();
    if (!newName || newName === oldName) {
      closeModal();
      return;
    }
    const parentPath = oldPath.substring(0, oldPath.lastIndexOf('/')) || '/';
    const newPath = (parentPath === '/' ? '/' : parentPath + '/') + newName;
    const res = await getHdfsSvc().Rename(oldPath, newPath);
    if (res.success) {
      showToast('重命名成功', 'success');
      closeModal();
      await refreshDir();
    } else {
      showToast(res.message, 'error');
    }
  });
}

// ===== Modal =====
function showModal(title, bodyHtml, onConfirm) {
  const overlay = document.getElementById('modal-overlay');
  overlay.innerHTML = `
    <div class="modal">
      <div class="modal-header">
        <span>${escapeHtml(title)}</span>
        <button class="close-btn" onclick="closeModal()">&times;</button>
      </div>
      <div class="modal-body">${bodyHtml}</div>
      <div class="modal-footer">
        <button class="btn" onclick="closeModal()">取消</button>
        <button class="btn btn-primary" id="modal-confirm-btn">确定</button>
      </div>
    </div>
  `;
  overlay.style.display = 'flex';
  document.getElementById('modal-confirm-btn').onclick = onConfirm;

  // Focus first input
  const firstInput = overlay.querySelector('input');
  if (firstInput) setTimeout(() => firstInput.focus(), 50);
}

function closeModal() {
  const overlay = document.getElementById('modal-overlay');
  overlay.style.display = 'none';
  overlay.innerHTML = '';
}

// ===== Context Menu =====
function showContextMenu(e, file) {
  e.preventDefault();
  hideContextMenu();
  const menu = document.getElementById('context-menu');
  const items = [];

  if (file) {
    if (file.isDir) {
      items.push({ label: '打开', action: () => navigateTo(file.path) });
      items.push({ divider: true });
    } else {
      items.push({ label: '下载', action: () => { state.selectedFile = file; downloadFile(); } });
      items.push({ divider: true });
    }
    items.push({ label: '重命名', action: () => { state.selectedFile = file; renameSelected(); } });
    items.push({ label: '删除', action: () => { state.selectedFile = file; deleteSelected(); }, danger: true });
  } else {
    items.push({ label: '新建文件夹', action: createFolder });
    items.push({ label: '上传文件', action: uploadFile });
    items.push({ divider: true });
    items.push({ label: '刷新', action: refreshDir });
  }

  menu.innerHTML = items.map(item => {
    if (item.divider) return '<div class="context-menu-divider"></div>';
    return `<div class="context-menu-item${item.danger ? ' danger' : ''}" data-action="${items.indexOf(item)}">${item.label}</div>`;
  }).join('');

  // Position
  menu.style.left = e.clientX + 'px';
  menu.style.top = e.clientY + 'px';
  menu.style.display = 'block';

  // Attach actions
  menu.querySelectorAll('.context-menu-item').forEach(el => {
    const idx = parseInt(el.dataset.action);
    el.onclick = () => { hideContextMenu(); items[idx].action(); };
  });
}

function hideContextMenu() {
  const menu = document.getElementById('context-menu');
  menu.style.display = 'none';
}

document.addEventListener('click', hideContextMenu);

// ===== Loading =====
function setLoading(val) {
  state.loading = val;
  const el = document.getElementById('loading-overlay');
  if (el) el.style.display = val ? 'flex' : 'none';
}

// ===== Utility =====
function escapeHtml(str) {
  const div = document.createElement('div');
  div.textContent = str;
  return div.innerHTML;
}

// ===== Render =====
function renderAll() {
  renderSidebar();
  renderContent();
  renderStatusBar();
  renderToolbar();
}

function renderSidebar() {
  const list = document.getElementById('sidebar-list');
  if (!state.connections || state.connections.length === 0) {
    list.innerHTML = '<div class="sidebar-empty">暂无连接<br>点击 + 添加HDFS连接</div>';
    return;
  }
  list.innerHTML = state.connections.map(conn => {
    const isActive = state.activeConnId === conn.id && state.connected;
    return `
      <div class="sidebar-item${isActive ? ' active connected' : ''}" data-id="${conn.id}">
        <span class="conn-icon"></span>
        <div class="conn-info" onclick="connectToHdfs('${conn.id}')">
          <div class="conn-name">${escapeHtml(conn.name)}</div>
          <div class="conn-addr">${escapeHtml(conn.address)}</div>
        </div>
        <div class="conn-actions">
          <button class="conn-action-btn" title="编辑" onclick="event.stopPropagation();editConnection('${conn.id}')">
            ${ICONS.edit}
          </button>
          <button class="conn-action-btn delete" title="删除" onclick="event.stopPropagation();deleteConnection('${conn.id}')">
            ${ICONS.trash}
          </button>
        </div>
      </div>
    `;
  }).join('');
}

function renderToolbar() {
  const hasSelection = !!state.selectedFile;
  const isFile = hasSelection && !state.selectedFile.isDir;
  document.getElementById('btn-upload').disabled = !state.connected;
  document.getElementById('btn-mkdir').disabled = !state.connected;
  document.getElementById('btn-download').disabled = !isFile;
  document.getElementById('btn-delete').disabled = !hasSelection;
  document.getElementById('btn-rename').disabled = !hasSelection;
  document.getElementById('btn-refresh').disabled = !state.connected;
  document.getElementById('btn-disconnect').disabled = !state.connected;
}

function renderContent() {
  const content = document.getElementById('file-content');
  const pathInput = document.getElementById('path-input');
  const btnBack = document.getElementById('btn-back');
  const btnUp = document.getElementById('btn-up');

  pathInput.value = state.currentPath;
  btnBack.disabled = state.historyIndex <= 0;
  btnUp.disabled = state.currentPath === '/';

  if (!state.connected) {
    document.getElementById('file-list-container').innerHTML = `
      <div class="welcome-state">
        ${ICONS.server}
        <div class="title">HDFS King</div>
        <div class="hint">从左侧选择一个连接开始浏览</div>
      </div>
    `;
    renderToolbar();
    return;
  }

  if (state.files.length === 0) {
    document.getElementById('file-list-container').innerHTML = `
      <div class="welcome-state">
        ${ICONS.folder}
        <div class="hint">空目录</div>
      </div>
    `;
    renderToolbar();
    return;
  }

  // Sort: folders first, then files
  const sorted = [...state.files].sort((a, b) => {
    if (a.isDir !== b.isDir) return a.isDir ? -1 : 1;
    return a.name.localeCompare(b.name);
  });

  const rows = sorted.map(f => {
    const selected = state.selectedFile && state.selectedFile.path === f.path;
    const iconClass = f.isDir ? 'folder' : 'file';
    const icon = f.isDir ? ICONS.folder : ICONS.file;
    return `
      <tr class="${selected ? 'selected' : ''}" data-path="${escapeHtml(f.path)}" data-isdir="${f.isDir}">
        <td>
          <div class="file-name-cell">
            <div class="file-icon ${iconClass}">${icon}</div>
            <span class="file-text-name">${escapeHtml(f.name)}</span>
          </div>
        </td>
        <td class="col-size">${f.isDir ? '-' : formatSize(f.size)}</td>
        <td class="col-time">${f.modTime || '-'}</td>
        <td class="col-owner">${escapeHtml(f.owner || '-')}</td>
        <td class="col-perm">${escapeHtml(f.permission || '-')}</td>
      </tr>
    `;
  }).join('');

  document.getElementById('file-list-container').innerHTML = `
    <table class="file-list-table">
      <thead>
        <tr>
          <th>名称</th>
          <th class="col-size">大小</th>
          <th class="col-time">修改时间</th>
          <th class="col-owner">所有者</th>
          <th class="col-perm">权限</th>
        </tr>
      </thead>
      <tbody id="file-tbody">${rows}</tbody>
    </table>
  `;

  // Attach events
  document.querySelectorAll('#file-tbody tr').forEach(row => {
    const filePath = row.dataset.path;
    const isDir = row.dataset.isdir === 'true';
    const file = state.files.find(f => f.path === filePath);

    row.onclick = (e) => {
      state.selectedFile = file;
      // Update selection visually without re-rendering the whole DOM
      document.querySelectorAll('#file-tbody tr').forEach(r => r.classList.remove('selected'));
      row.classList.add('selected');
      renderToolbar();
    };

    row.ondblclick = () => {
      if (isDir) navigateTo(filePath);
    };

    row.oncontextmenu = (e) => {
      state.selectedFile = file;
      document.querySelectorAll('#file-tbody tr').forEach(r => r.classList.remove('selected'));
      row.classList.add('selected');
      renderToolbar();
      showContextMenu(e, file);
    };
  });

  // Right click on empty area
  document.getElementById('file-list-container').oncontextmenu = (e) => {
    if (e.target.closest('#file-tbody tr')) return;
    showContextMenu(e, null);
  };

  renderToolbar();
}

function renderStatusBar() {
  const connStatus = document.getElementById('status-conn-text');
  const dot = document.getElementById('status-dot');
  const fileCount = document.getElementById('status-count');

  if (state.connected) {
    const conn = state.connections.find(c => c.id === state.activeConnId);
    connStatus.textContent = conn ? `已连接: ${conn.name}` : '已连接';
    dot.classList.add('connected');
  } else {
    connStatus.textContent = '未连接';
    dot.classList.remove('connected');
  }

  const dirs = state.files.filter(f => f.isDir).length;
  const files = state.files.filter(f => !f.isDir).length;
  fileCount.textContent = state.connected ? `${dirs} 个文件夹, ${files} 个文件` : '';
}

// ===== Init DOM =====
document.querySelector('#app').innerHTML = `
  <!-- Toolbar -->
  <div class="toolbar">
    <div class="toolbar-group">
      <button class="toolbar-btn" id="btn-upload" disabled title="上传文件">
        ${ICONS.upload}<span>上传</span>
      </button>
      <button class="toolbar-btn" id="btn-download" disabled title="下载文件">
        ${ICONS.download}<span>下载</span>
      </button>
    </div>
    <div class="toolbar-divider"></div>
    <div class="toolbar-group">
      <button class="toolbar-btn" id="btn-mkdir" disabled title="新建文件夹">
        ${ICONS.folderPlus}<span>新建文件夹</span>
      </button>
      <button class="toolbar-btn" id="btn-rename" disabled title="重命名">
        ${ICONS.edit}<span>重命名</span>
      </button>
      <button class="toolbar-btn" id="btn-delete" disabled title="删除">
        ${ICONS.trash}<span>删除</span>
      </button>
    </div>
    <div class="toolbar-divider"></div>
    <div class="toolbar-group">
      <button class="toolbar-btn" id="btn-refresh" disabled title="刷新">
        ${ICONS.refresh}<span>刷新</span>
      </button>
    </div>
    <div style="flex:1"></div>
    <div class="toolbar-group">
      <button class="toolbar-btn" id="btn-disconnect" disabled title="断开连接">
        ${ICONS.disconnect}<span>断开</span>
      </button>
    </div>
  </div>

  <!-- Main -->
  <div class="main-container">
    <!-- Sidebar -->
    <div class="sidebar">
      <div class="sidebar-header">
        <span>连接管理</span>
        <button title="添加连接" id="btn-add-conn">${ICONS.plus}</button>
      </div>
      <div class="sidebar-list" id="sidebar-list"></div>
    </div>

    <!-- Content -->
    <div class="content">
      <div class="breadcrumb-bar">
        <button class="nav-btn" id="btn-back" title="后退" disabled>${ICONS.back}</button>
        <button class="nav-btn" id="btn-up" title="上级目录" disabled>${ICONS.up}</button>
        <button class="nav-btn" id="btn-home" title="根目录">${ICONS.home}</button>
        <input class="path-input" id="path-input" value="/" />
      </div>
      <div class="content-wrapper" style="position:relative;flex:1;display:flex;flex-direction:column;overflow:hidden;">
        <div class="file-list-container" id="file-list-container">
          <div class="welcome-state">
            ${ICONS.server}
            <div class="title">HDFS King</div>
            <div class="hint">从左侧选择一个连接开始浏览</div>
          </div>
        </div>
        <div class="loading-overlay" id="loading-overlay" style="display:none;">
          <div class="spinner"></div>
        </div>
      </div>
    </div>
  </div>

  <!-- Status Bar -->
  <div class="status-bar">
    <div class="status-conn">
      <span class="status-dot" id="status-dot"></span>
      <span id="status-conn-text">未连接</span>
    </div>
    <span id="status-count"></span>
  </div>

  <!-- Modal -->
  <div class="modal-overlay" id="modal-overlay" style="display:none;"></div>

  <!-- Context Menu -->
  <div class="context-menu" id="context-menu" style="display:none;"></div>

  <!-- Toast -->
  <div class="toast-container" id="toast-container"></div>
`;

// ===== Bind Events =====
document.getElementById('btn-add-conn').onclick = addConnection;
document.getElementById('btn-upload').onclick = uploadFile;
document.getElementById('btn-download').onclick = downloadFile;
document.getElementById('btn-mkdir').onclick = createFolder;
document.getElementById('btn-rename').onclick = renameSelected;
document.getElementById('btn-delete').onclick = deleteSelected;
document.getElementById('btn-refresh').onclick = refreshDir;
document.getElementById('btn-disconnect').onclick = disconnectHdfs;
document.getElementById('btn-back').onclick = goBack;
document.getElementById('btn-up').onclick = goUp;
document.getElementById('btn-home').onclick = () => { if (state.connected) navigateTo('/'); };

// Path input: enter to navigate
document.getElementById('path-input').onkeydown = (e) => {
  if (e.key === 'Enter') {
    const path = e.target.value.trim();
    if (path && state.connected) navigateTo(path);
  }
};

// Keyboard shortcuts
document.addEventListener('keydown', (e) => {
  if (document.getElementById('modal-overlay').style.display === 'flex') return;
  if (e.key === 'Delete' && state.selectedFile) {
    deleteSelected();
  } else if (e.key === 'F2' && state.selectedFile) {
    renameSelected();
  } else if (e.key === 'F5') {
    e.preventDefault();
    refreshDir();
  } else if (e.key === 'Backspace' && !e.target.matches('input')) {
    goUp();
  }
});

// Expose functions to global (for inline onclick handlers)
window.connectToHdfs = connectToHdfs;
window.editConnection = editConnection;
window.deleteConnection = deleteConnection;
window.closeModal = closeModal;
window.navigateTo = navigateTo;

// ===== Initial load =====
loadConnections();
