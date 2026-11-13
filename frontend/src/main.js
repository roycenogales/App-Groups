import { WindowService, WorkspaceService } from "../bindings/changeme";

const scanButton = document.getElementById("scan-button");
const selectAllButton = document.getElementById("select-all-button");
const clearSelectionButton = document.getElementById("clear-selection-button");
const saveWorkspaceButton = document.getElementById("save-workspace-button");
const workspaceNameInput = document.getElementById("workspace-name");

const windowList = document.getElementById("window-list");
const workspaceList = document.getElementById("workspace-list");
const windowCount = document.getElementById("window-count");
const status = document.getElementById("status");

let windows = [];
let selectedHandles = new Set();
let workspaces = [];

function updateSelectionCount() {
    const selected = selectedHandles.size;
    const total = windows.length;

    windowCount.textContent =
        `${total} detected · ${selected} selected`;

    selectAllButton.disabled = total === 0;
    clearSelectionButton.disabled = total === 0;
    saveWorkspaceButton.disabled = selected === 0;
}

function renderWindows() {
    windowList.replaceChildren();

    if (windows.length === 0) {
        const empty = document.createElement("div");
        empty.className = "empty-state";
        empty.innerHTML = `
            <div class="empty-icon">▣</div>
            <h3>No windows found</h3>
            <p>No eligible visible desktop windows were detected.</p>
        `;
        windowList.appendChild(empty);
        updateSelectionCount();
        return;
    }

    for (const win of windows) {
        const card = document.createElement("article");
        card.className = "window-card";

        const header = document.createElement("div");
        header.className = "window-card-header";

        const checkbox = document.createElement("input");
        checkbox.type = "checkbox";
        checkbox.className = "window-checkbox";
        checkbox.checked = selectedHandles.has(win.hwnd);
        checkbox.setAttribute(
            "aria-label",
            `Include ${win.process_name} in workspace`
        );

        checkbox.addEventListener("change", () => {
            if (checkbox.checked) {
                selectedHandles.add(win.hwnd);
            } else {
                selectedHandles.delete(win.hwnd);
            }

            updateSelectionCount();
        });

        const identity = document.createElement("div");
        identity.className = "window-identity";

        const processName = document.createElement("div");
        processName.className = "process-name";
        processName.textContent = win.process_name || "Unknown process";

        const title = document.createElement("div");
        title.className = "window-title";
        title.textContent = win.title || "Untitled window";

        identity.append(processName, title);
        header.append(checkbox, identity);

        const details = document.createElement("div");
        details.className = "window-details";

        const geometry = document.createElement("span");
        geometry.textContent =
            `${win.width} × ${win.height} at (${win.x}, ${win.y})`;

        const state = document.createElement("span");
        state.className = `state state-${win.state}`;
        state.textContent = win.state;

        details.append(geometry, state);

        const path = document.createElement("div");
        path.className = "process-path";
        path.textContent = win.process_path || `PID: ${win.pid}`;
        path.title = win.process_path || "";

        card.append(header, details, path);
        windowList.appendChild(card);
    }

    updateSelectionCount();
}

function renderWorkspaces() {
    workspaceList.replaceChildren();

    if (workspaces.length === 0) {
        const empty = document.createElement("div");
        empty.className = "empty-workspaces";
        empty.textContent = "No saved workspaces yet.";
        workspaceList.appendChild(empty);
        return;
    }

    for (const workspace of workspaces) {
        const item = document.createElement("div");
        item.className = "workspace-item";

        const info = document.createElement("div");
        info.className = "workspace-info";

        const name = document.createElement("div");
        name.className = "workspace-name";
        name.textContent = workspace.name;

        const count = document.createElement("div");
        count.className = "workspace-count";
        const windowTotal = workspace.windows?.length ?? 0;
        count.textContent =
            `${windowTotal} window${windowTotal === 1 ? "" : "s"}`;

        info.append(name, count);

        const deleteButton = document.createElement("button");
        deleteButton.className = "delete-workspace-button";
        deleteButton.textContent = "×";
        deleteButton.title = `Delete ${workspace.name}`;
        deleteButton.setAttribute(
            "aria-label",
            `Delete workspace ${workspace.name}`
        );

        deleteButton.addEventListener("click", () => {
            deleteWorkspace(workspace);
        });

        const actions = document.createElement("div");
        actions.className = "workspace-actions";

        const restoreButton = document.createElement("button");
        restoreButton.className = "restore-workspace-button";
        restoreButton.textContent = "Restore";
        restoreButton.title = `Restore ${workspace.name}`;

        restoreButton.addEventListener("click", () => {
            restoreWorkspace(workspace);
        });

        actions.append(restoreButton, deleteButton);

        item.append(info, actions);
        workspaceList.appendChild(item);
    }
}

async function loadWorkspaces() {
    try {
        workspaces = await WorkspaceService.ListWorkspaces();
        renderWorkspaces();
    } catch (error) {
        console.error("Unable to load workspaces:", error);

        workspaceList.textContent = "Unable to load saved workspaces.";
        status.textContent =
            `Unable to load workspaces: ${error?.message ?? error}`;
    }
}

async function scanWindows() {
    scanButton.disabled = true;
    scanButton.textContent = "Scanning...";
    status.textContent = "Scanning visible desktop windows...";

    try {
        windows = await WindowService.ListWindows();

        // All detected windows are included by default.
        selectedHandles = new Set(windows.map(win => win.hwnd));

        renderWindows();

        status.textContent =
            `Scan complete. Found ${windows.length} window(s).`;
    } catch (error) {
        console.error("Window scan failed:", error);
        status.textContent =
            `Unable to scan windows: ${error?.message ?? error}`;

        windowCount.textContent = "Scan failed";
    } finally {
        scanButton.disabled = false;
        scanButton.innerHTML = "<span aria-hidden=\"true\">↻</span> Scan Windows";
    }
}

async function saveWorkspace() {
    const name = workspaceNameInput.value.trim();

    if (!name) {
        status.textContent = "Enter a name for this workspace.";
        workspaceNameInput.focus();
        return;
    }

    const selectedWindows = windows.filter(
        win => selectedHandles.has(win.hwnd)
    );

    if (selectedWindows.length === 0) {
        status.textContent = "Select at least one window before saving.";
        return;
    }

    saveWorkspaceButton.disabled = true;
    saveWorkspaceButton.textContent = "Saving...";

    try {
        const saved = await WorkspaceService.SaveWorkspace(
            name,
            selectedWindows
        );

        workspaces.push(saved);
        renderWorkspaces();

        workspaceNameInput.value = "";

        status.textContent =
            `Workspace "${saved.name}" saved with ${saved.windows.length} window(s).`;
    } catch (error) {
        console.error("Unable to save workspace:", error);
        status.textContent =
            `Unable to save workspace: ${error?.message ?? error}`;
    } finally {
        saveWorkspaceButton.textContent = "Save Workspace";
        updateSelectionCount();
    }
}

async function restoreWorkspace(workspace) {
    status.textContent =
        `Restoring workspace "${workspace.name}"...`;

    try {
        await WorkspaceService.RestoreWorkspace(workspace.id);

        status.textContent =
            `Workspace "${workspace.name}" restored successfully.`;
    } catch (error) {
        console.error("Unable to restore workspace:", error);

        status.textContent =
            `Unable to restore workspace: ${error?.message ?? error}`;
    }
}

async function deleteWorkspace(workspace) {
    const confirmed = window.confirm(
        `Delete the workspace "${workspace.name}"?`
    );

    if (!confirmed) {
        return;
    }

    try {
        await WorkspaceService.DeleteWorkspace(workspace.id);

        workspaces = workspaces.filter(
            item => item.id !== workspace.id
        );

        renderWorkspaces();

        status.textContent =
            `Workspace "${workspace.name}" deleted.`;
    } catch (error) {
        console.error("Unable to delete workspace:", error);
        status.textContent =
            `Unable to delete workspace: ${error?.message ?? error}`;
    }
}

scanButton.addEventListener("click", scanWindows);

selectAllButton.addEventListener("click", () => {
    selectedHandles = new Set(windows.map(win => win.hwnd));
    renderWindows();
});

clearSelectionButton.addEventListener("click", () => {
    selectedHandles.clear();
    renderWindows();
});

saveWorkspaceButton.addEventListener("click", saveWorkspace);

workspaceNameInput.addEventListener("keydown", event => {
    if (event.key === "Enter" && !saveWorkspaceButton.disabled) {
        saveWorkspace();
    }
});

// Load saved workspaces as soon as the frontend starts.
loadWorkspaces();
