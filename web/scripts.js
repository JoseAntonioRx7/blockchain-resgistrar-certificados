const API_URL = "http://localhost:8080";


// -------------------------------
// 1. INICIALIZAÇÃO
// -------------------------------

document.addEventListener("DOMContentLoaded", () => {
    checkAuth();
    updateCurrentDate();
    loadDashboard();

    const forms = {
        loginForm: loginInstituicao,
        registerForm: registerCertificate,
        adminForm: registerInstitution
    };

    Object.keys(forms).forEach(id => {
        const el = document.getElementById(id);
        if (el) el.addEventListener("submit", forms[id]);
    });
});


// -------------------------------
// 2. UTILITÁRIOS
// -------------------------------

function toggleDisplay(id, value) {
    const el = document.getElementById(id);
    if (el) el.style.display = value;
}

function setLoading(btn, isLoading, text) {
    btn.disabled = isLoading;
    btn.innerHTML = text;
}

function updateCurrentDate() {
    const dateElement = document.getElementById("currentDate");
    if (!dateElement) return;

    const options = {
        weekday: "long",
        year: "numeric",
        month: "long",
        day: "numeric"
    };

    dateElement.innerText = new Date().toLocaleDateString("pt-BR", options);
}


// -------------------------------
// 3. AUTENTICAÇÃO
// -------------------------------

function checkAuth() {
    const token = localStorage.getItem("instituicao_token");

    const ui = {
        login: ["loginSection"],
        dashboard: ["dashboardStats", "mainActions", "authHeader"]
    };

    // Admin menu e seção visíveis apenas quando logado
    const adminElements = ["adminMenuItem", "adminSection"];
    const isLoggedIn = !!token;

    adminElements.forEach(id => {
        toggleDisplay(id, isLoggedIn ? "block" : "none");
    });

    if (token) {
        ui.login.forEach(id => toggleDisplay(id, "none"));
        ui.dashboard.forEach(id =>
            toggleDisplay(id, id === "authHeader" ? "flex" : "block")
        );

        document.getElementById("pageTitle").innerText = "Visão Geral do Sistema";
    } else {
        ui.login.forEach(id => toggleDisplay(id, "block"));
        ui.dashboard.forEach(id => toggleDisplay(id, "none"));

        document.getElementById("pageTitle").innerText = "Acesso Restrito";
    }
}

function logoutInstituicao() {
    localStorage.removeItem("instituicao_token");
    checkAuth();
}


// -------------------------------
// 4. LOGIN
// -------------------------------

async function loginInstituicao(e) {
    e.preventDefault();

    const user = document.getElementById("loginUser").value;
    const pass = document.getElementById("loginPass").value;
    const resultDiv = document.getElementById("loginResult");

    try {
        const response = await fetch(`${API_URL}/login`, {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({
                username: user,
                password: pass
            })
        });

        const data = await response.json();

        if (!response.ok) {
            resultDiv.innerText = data.error || "Credenciais inválidas";
            resultDiv.style.color = "#ef4444";
            return;
        }

        localStorage.setItem("instituicao_token", data.token);

        resultDiv.innerHTML =
            `<span style="color:#10b981">Entrando na rede...</span>`;

        setTimeout(() => {
            checkAuth();
            loadDashboard();
        }, 800);

    } catch (err) {
        resultDiv.innerText = "Erro na conexão com a rede blockchain.";
    }
}


// -------------------------------
// 5. REGISTRAR CERTIFICADO
// -------------------------------

async function registerCertificate(e) {
    e.preventDefault();

    const token = localStorage.getItem("instituicao_token");
    const btn = e.target.querySelector("button");
    const resultDiv = document.getElementById("registerResult");

    if (!token) {
        resultDiv.innerHTML = `<span style="color:#ef4444">Token não encontrado. Faça login novamente.</span>`;
        return;
    }

    setLoading(btn, true, "Minerando bloco...");

    try {
        // Cria o FormData a partir do formulário
        const formData = new FormData(e.target);
        
        // Adiciona o token ao FormData
        formData.append("token", token);

        console.log("Enviando certificado com token:", token.substring(0, 20) + "...");

        const response = await fetch(`${API_URL}/register`, {
            method: "POST",
            headers: {
                // Enviamos o token tanto no header quanto no body
                "Authorization": "Bearer " + token
            },
            body: formData
        });

        console.log("Response status:", response.status);

        const data = await response.json();
        console.log("Response data:", data);

        if (!response.ok) {
            resultDiv.innerHTML =
                `<span style="color:#ef4444">${data.error || "Erro desconhecido"}</span>`;
            return;
        }

        resultDiv.innerHTML =
            `<div class="alert-success">
            Bloco minerado! ID: ${data.id ? data.id.substring(0,8) : 'N/A'}
            </div>`;

        e.target.reset();
        loadDashboard();

    } catch (err) {
        console.error("Erro completo:", err);
        resultDiv.innerText = "Erro no processo de registro: " + err.message;
    }

    setLoading(btn, false, "Garantir Autenticidade");
}


// -------------------------------
// 6. REGISTRAR INSTITUIÇÃO
// -------------------------------

async function registerInstitution(e) {
    e.preventDefault();

    const resultDiv = document.getElementById("adminResult");
    const btn = e.target.querySelector("button");

    const body = {
        name: document.getElementById("instName").value,
        username: document.getElementById("instUser").value,
        password: document.getElementById("instPass").value
    };

    setLoading(btn, true, "Gerando credenciais...");

    try {
        const response = await fetch(
            `${API_URL}/admin/register-institution`,
            {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify(body)
            }
        );

        const data = await response.json();

        if (!response.ok) {
            resultDiv.innerHTML =
                `<span style="color:#ef4444">${data.error}</span>`;
            return;
        }

        resultDiv.innerHTML =
            `<div class="alert-info">
            Chave privada: <code>${data.private_key}</code>
            </div>`;

        e.target.reset();

    } catch (err) {
        resultDiv.innerText = "Erro ao registrar instituição.";
    }

    setLoading(btn, false, "Criar Credenciais");
}


// -------------------------------
// 7. DASHBOARD
// -------------------------------

async function loadDashboard() {

    const token = localStorage.getItem("instituicao_token");
    const tbody = document.getElementById("dashboardBody");

    if (!token || !tbody) return;

    try {
        const response = await fetch(`${API_URL}/list`, {
            headers: {
                Authorization: `Bearer ${token}`
            }
        });

        if (!response.ok) {
            throw new Error("Falha ao carregar dashboard: " + response.status);
        }

        const data = await response.json();
        
        // Handle case where response might be an error object
        if (data.error) {
            console.error("Erro da API:", data.error);
            return;
        }

        const certs = Array.isArray(data) ? data : [];
        
        // Ensure we have an array before accessing .length
        if (!certs || !Array.isArray(certs)) {
            console.error("Dados inesperados:", data);
            return;
        }

        document.getElementById("statTotal").innerText = certs.length;
        document.getElementById("statBlocks").innerText = certs.length;
        document.getElementById("countCerts").innerText =
            `${certs.length} certificados`;

        if (certs.length === 0) {
            tbody.innerHTML = '<tr><td colspan="5" style="text-align:center">Nenhum certificado encontrado</td></tr>';
            return;
        }

        tbody.innerHTML = certs.reverse().map(cert => `
            <tr>
                <td><strong>${cert.student_name}</strong></td>
                <td>${cert.course}<br><small>${cert.institution || ''}</small></td>
                <td>${cert.timestamp ? new Date(cert.timestamp * 1000).toLocaleDateString() : 'N/A'}</td>
                <td><code>${(cert.file_hash || '').substring(0,8)}...</code></td>
                <td style="text-align:right">
                    <button onclick="copyHash('${cert.file_hash}')">📋</button>
                    <a href="${API_URL}/pdfs/cert_${cert.id}.pdf" target="_blank">PDF</a>
                </td>
            </tr>
        `).join("");

    } catch (err) {
        console.error("Erro no dashboard:", err);
    }
}


// -------------------------------
// 8. VERIFICAÇÃO DE CERTIFICADO
// -------------------------------

async function verifyCertificate() {

    const hash = document.getElementById("hashInput").value.trim();
    const resultDiv = document.getElementById("verifyResult");

    if (!hash) return;

    resultDiv.innerText = "Consultando blockchain...";

    try {
        const response = await fetch(
            `${API_URL}/verify?hash=${encodeURIComponent(hash)}`
        );

        const data = await response.json();

        if (!response.ok || data.error) {
            resultDiv.innerHTML =
                `<div class="alert-error">Documento não encontrado</div>`;
            return;
        }

        resultDiv.innerHTML = `
            <div class="result-card">
                <h4>Certificado Autêntico</h4>
                <p><strong>Aluno:</strong> ${data.student_name}</p>
                <p><strong>Hash:</strong> ${hash.substring(0,12)}...</p>
            </div>
        `;

    } catch (err) {
        resultDiv.innerText = "Erro na rede de validação.";
    }
}


// -------------------------------
// 9. FILTRO
// -------------------------------

function filterTable() {
    const filter = document
        .getElementById("searchInput")
        .value
        .toUpperCase();

    document
        .querySelectorAll("#dashboardBody tr")
        .forEach(row => {
            row.style.display =
                row.innerText.toUpperCase().includes(filter)
                    ? ""
                    : "none";
        });
}

function copyHash(hash) {
    navigator.clipboard.writeText(hash);
    alert("Hash copiado!");
}
