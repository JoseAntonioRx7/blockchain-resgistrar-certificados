/**
 * TTLedger - Blockchain Infrastructure
 * Versão Final Consolidada e Corrigida
 */

const BASE_URL = "http://localhost:8080/api";

// ---------------------------------------------------------
// 1. UTILITÁRIOS DE INTERFACE (UI)
// ---------------------------------------------------------
const UI = {
    toggle: (id, display) => {
        const el = document.getElementById(id);
        if (el) el.style.display = display;
    },
    text: (id, txt) => {
        const el = document.getElementById(id);
        if (el) el.innerText = txt;
    },
    setLoading: (btn, isLoading, text) => {
        if (btn) {
            btn.disabled = isLoading;
            btn.innerHTML = isLoading ? '<i class="fas fa-spinner fa-spin"></i>' : text;
        }
    },
    showResult: (id, message, isError = false) => {
        const el = document.getElementById(id);
        if (el) {
            el.innerHTML = message;
            el.className = `result ${isError ? 'error-text' : 'success-text'}`;
            el.style.display = "block";
            if (!isError && (message.includes('code') || message.includes('HEX'))) {
                el.style.background = "rgba(0, 255, 136, 0.1)";
                el.style.border = "1px solid #00ff88";
            }
        }
    }
};

// ---------------------------------------------------------
// 2. COMUNICAÇÃO COM API
// ---------------------------------------------------------
async function safeFetch(endpoint, options = {}) {
    const cleanEndpoint = endpoint.startsWith('/') ? endpoint.slice(1) : endpoint;
    const url = `${BASE_URL}/${cleanEndpoint}`;
    
    try {
        const response = await fetch(url, options);
        if (response.status === 401) {
            localStorage.clear();
            location.reload();
            throw new Error("Sessão expirada. Refaça o login.");
        }
        
        const data = await response.json();
        if (!response.ok) throw new Error(data.error || "Erro na requisição");
        return data;
    } catch (err) {
        console.error("Erro API:", err.message);
        throw err;
    }
}

// ---------------------------------------------------------
// 3. NAVEGAÇÃO E CONTROLE DE ACESSO
// ---------------------------------------------------------
function showTab(tabId) {
    UI.toggle("dashboardStats", "none");
    UI.toggle("mainActions", "none");
    UI.toggle("adminSection", "none");
    UI.toggle("loginSection", "none");

    if (tabId === 'dashboard') {
        UI.toggle("dashboardStats", "grid");
        UI.toggle("mainActions", "grid");
        UI.text("pageTitle", "Dashboard Principal");
    } else if (tabId === 'admin') {
        UI.toggle("adminSection", "block");
        UI.text("pageTitle", "Expansão de Rede");
    }

    document.querySelectorAll(".sidebar-nav li").forEach(li => li.classList.remove("active"));
    const activeLink = document.querySelector(`.sidebar-nav a[href*="${tabId}"]`);
    if (activeLink) activeLink.parentElement.classList.add("active");
}

function checkAuth() {
    const token = localStorage.getItem("ttledger_token");
    const user = localStorage.getItem("ttledger_user");
    const isLoggedIn = !!token;

    UI.toggle("dashboardContent", isLoggedIn ? "block" : "none");
    UI.toggle("loginSection", isLoggedIn ? "none" : "block");
    UI.toggle("authHeader", isLoggedIn ? "flex" : "none");
    UI.text("userNameDisplay", user || "Instituição");

    if (isLoggedIn) {
        const adminMenu = document.getElementById("adminMenuItem");
        if (adminMenu) {
            adminMenu.style.display = (user === 'admin') ? "block" : "none";
        }
        showTab('dashboard');
        loadDashboard();
    }
}

// ---------------------------------------------------------
// 4. HANDLERS (Ações Globais)
// ---------------------------------------------------------

/**
 * RESOLUÇÃO: Função de Verificação (Exposta globalmente)
 */
async function verifyCertificate() {
    const hash = document.getElementById("hashInput").value.trim();
    if (!hash) {
        UI.showResult("verifyResult", "Por favor, insira um hash para verificar.", true);
        return;
    }

    try {
        const data = await safeFetch(`verify?hash=${hash}`);
        UI.showResult("verifyResult", `
            <div style="text-align: left; border-left: 3px solid #00ff88; padding-left: 15px;">
                <p>✅ <strong>Certificado Autêntico</strong></p>
                <p><strong>Aluno:</strong> ${data.student_name}</p>
                <p><strong>Curso:</strong> ${data.course}</p>
                <p><strong>Instituição:</strong> ${data.institution}</p>
                <p><strong>Data:</strong> ${new Date(data.timestamp * 1000).toLocaleString()}</p>
            </div>
        `, false);
    } catch (err) {
        UI.showResult("verifyResult", "Certificado não encontrado ou inválido.", true);
    }
}

async function handleLogin(e) {
    e.preventDefault();
    const btn = e.target.querySelector("button");
    UI.setLoading(btn, true, "Autenticando...");

    try {
        const username = e.target.loginUser.value;
        const data = await safeFetch("login", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ username, password: e.target.loginPass.value })
        });

        localStorage.setItem("ttledger_token", data.token);
        localStorage.setItem("ttledger_user", username);
        location.reload();
    } catch (err) {
        UI.showResult("loginResult", err.message, true);
        UI.setLoading(btn, false, "Entrar no Sistema");
    }
}

async function handleRegisterCertificate(e) {
    e.preventDefault();
    const btn = e.target.querySelector("button");
    const token = localStorage.getItem("ttledger_token");
    UI.setLoading(btn, true, "Minerando na Rede...");

    try {
        const formData = new FormData(e.target);
        const response = await fetch(`${BASE_URL}/register`, {
            method: "POST",
            headers: { "Authorization": `Bearer ${token}` },
            body: formData
        });

        const data = await response.json();
        if (!response.ok) throw new Error(data.error || "Erro no registro");

        UI.showResult("registerResult", "✅ Certificado registrado com sucesso!");
        e.target.reset();
        loadDashboard();
    } catch (err) {
        UI.showResult("registerResult", err.message, true);
    } finally {
        UI.setLoading(btn, false, "Garantir Autenticidade");
    }
}

async function loadDashboard() {
    const token = localStorage.getItem("ttledger_token");
    try {
        const data = await safeFetch("list", {
            headers: { "Authorization": `Bearer ${token}` }
        });
        renderTable(data);
    } catch (err) { console.error("Erro ao carregar lista:", err); }
}

/**
 * RESOLUÇÃO: Tabela organizada com links de visualização
 */
function renderTable(certs) {
    const tbody = document.getElementById("dashboardBody");
    if (!tbody) return;
    
    const list = Array.isArray(certs) ? certs : [];
    UI.text("countCerts", `${list.length} certificados`);
    UI.text("statTotal", list.length);

    tbody.innerHTML = list.length === 0 
        ? '<tr><td colspan="5">Nenhum registro encontrado nesta rede.</td></tr>'
        : list.map(c => `
            <tr>
                <td><strong>${c.student_name}</strong></td>
                <td>${c.course}</td>
                <td>${new Date(c.timestamp * 1000).toLocaleDateString()}</td>
                <td><code title="${c.file_hash}">${c.file_hash.substring(0, 12)}...</code></td>
                <td style="text-align: right;">
                <a href="/api/pdfs/cert_${c.file_hash}.pdf" target="_blank" class="btn-view" title="Visualizar PDF">
                    <i class="fas fa-file-pdf"></i>
                </a>
                    <button class="btn-view" onclick="navigator.clipboard.writeText('${c.file_hash}'); alert('Hash copiado!')" title="Copiar Hash">
                        <i class="fas fa-copy"></i>
                    </button>
                </td>
            </tr>
        `).join("");
}

// ---------------------------------------------------------
// 5. INICIALIZAÇÃO E EVENTOS
// ---------------------------------------------------------
document.addEventListener("DOMContentLoaded", () => {
    checkAuth();
    UI.text("currentDate", new Date().toLocaleDateString('pt-BR', { dateStyle: 'full' }));

    // Listeners de Formulários
    document.getElementById("loginForm")?.addEventListener("submit", handleLogin);
    document.getElementById("registerForm")?.addEventListener("submit", handleRegisterCertificate);
    
    // Global Logout
    window.logout = () => {
        localStorage.clear();
        location.reload();
    };

    // Vincula a função de verificação ao objeto window para o HTML acessar
    window.verifyCertificate = verifyCertificate;
});