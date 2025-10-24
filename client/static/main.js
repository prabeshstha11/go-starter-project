const API_BASE = "http://localhost:8800";

function loadTodos() {
    fetch(`${API_BASE}/todo`)
        .then((res) => res.json())
        .then((data) => {
            const container = document.getElementById("todo-list");
            container.innerHTML = "";
            let total = 0,
                pending = 0,
                completed = 0;

            data.todos.forEach((todo) => {
                total++;
                todo.isCompleted ? completed++ : pending++;

                const todoEl = document.createElement("div");
                todoEl.className = "flex justify-between my-3 p-3 border-1 todo-item";
                todoEl.innerHTML = `
                            <div class="flex items-center gap-3">
                                <input type="checkbox" ${todo.isCompleted ? "checked" : ""} />
                                <span class="text-3xl text-primary ${todo.isCompleted ? "line-through text-gray-400" : ""}">
                                    ${todo.item}
                                </span>
                            </div>
                            <div>
                                <button class="btn btn-success">Edit</button>
                                <button class="btn btn-secondary">Delete</button>
                            </div>
                        `;

                todoEl.querySelector("input[type=checkbox]").addEventListener("change", (e) => {
                    fetch(`${API_BASE}/todo/${todo.id}`, {
                        method: "PATCH",
                        headers: { "Content-Type": "application/json" },
                        body: JSON.stringify({ isCompleted: e.target.checked }),
                    }).then(loadTodos);
                });

                todoEl.querySelector(".btn-secondary").addEventListener("click", () => {
                    if (!confirm("Are you sure?")) return;
                    fetch(`${API_BASE}/todo/${todo.id}`, { method: "DELETE" }).then(loadTodos);
                });

                todoEl.querySelector(".btn-success").addEventListener("click", () => {
                    const input = document.getElementById("task-input");
                    input.value = todo.item;
                    const addBtn = document.getElementById("add-btn");
                    addBtn.textContent = "Edit";
                    addBtn.onclick = () => {
                        fetch(`${API_BASE}/todo/${todo.id}`, {
                            method: "PATCH",
                            headers: { "Content-Type": "application/json" },
                            body: JSON.stringify({ item: input.value }),
                        }).then(() => {
                            input.value = "";
                            addBtn.textContent = "Add Task";
                            addBtn.onclick = addTodo;
                            loadTodos();
                        });
                    };
                });

                container.appendChild(todoEl);
            });

            document.getElementById("total").textContent = total;
            document.getElementById("pending").textContent = pending;
            document.getElementById("completed").textContent = completed;
        });
}

function addTodo() {
    const input = document.getElementById("task-input");
    const value = input.value.trim();
    if (!value) return;

    fetch(`${API_BASE}/create`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ item: value, isCompleted: false }),
    }).then(() => {
        input.value = "";
        loadTodos();
    });
}

document.getElementById("add-btn").onclick = addTodo;
loadTodos();
