// Função para Registrar
document.getElementById('registerForm').addEventListener('submit', async (e) => {
    e.preventDefault();
    const resultDiv = document.getElementById('registerResult');
    resultDiv.innerText = "Minerando bloco... aguarde.";

    const formData = new FormData(e.target);

    try {
        const response = await fetch('/register', {
            method: 'POST',
            body: formData
        });

        const data = await response.json();
        if (response.ok) {
            resultDiv.innerHTML = `<strong>Sucesso!</strong><br>ID: ${data.id}<br>Hash: ${data.hash}`;
        } else {
            resultDiv.innerText = "Erro: " + (data.error || "Falha no registro");
        }
    } catch (err) {
        resultDiv.innerText = "Erro ao conectar com o servidor.";
    }
});

// Função para Verificar
async function verifyCertificate() {
    const hash = document.getElementById('hashInput').value;
    const resultDiv = document.getElementById('verifyResult');
    
    if (!hash) {
        resultDiv.innerText = "Por favor, insira um hash.";
        return;
    }

    try {
        // Corresponde ao r.URL.Query().Get("hash") do seu Go
        const response = await fetch(`/verify?hash=${hash}`);
        const data = await response.json();

        if (response.ok && !data.error) {
            resultDiv.innerHTML = `
                <strong>Certificado Válido!</strong><br>
                Aluno: ${data.student_name}<br>
                Instituição: ${data.institution}<br>
                Data: ${new Date(data.timestamp * 1000).toLocaleDateString()}
            `;
        } else {
            resultDiv.innerText = "Certificado não encontrado ou inválido.";
        }
    } catch (err) {
        resultDiv.innerText = "Erro ao consultar servidor.";
    }
}

async function loadDashboard() {
    const tbody = document.getElementById('dashboardBody');
    
    try {
        const response = await fetch('/list');
        const certs = await response.json();

        tbody.innerHTML = ""; // Limpa a tabela

        certs.forEach(cert => {
            const row = `
                <tr>
                    <td>${cert.student_name}</td>
                    <td>${cert.course}</td>
                    <td style="font-size: 10px;">${cert.file_hash}</td>
                </tr>
            `;
            tbody.innerHTML += row;
        });
    } catch (err) {
        console.error("Erro ao carregar dashboard:", err);
    }
}

// Carrega automaticamente ao abrir a página
window.onload = loadDashboard;


// Função para buscar e carregar os dados (já tínhamos começado)
async function loadDashboard() {
    const tbody = document.getElementById('dashboardBody');
    try {
        const response = await fetch('/list');
        const certs = await response.json();

        tbody.innerHTML = ""; 

        if (certs && certs.length > 0) {
            certs.forEach(cert => {
                const date = new Date(cert.timestamp * 1000).toLocaleDateString();
                const row = `
                    <tr style="border-bottom: 1px solid #eee;">
                        <td style="padding: 12px;"><strong>${cert.student_name}</strong></td>
                        <td style="padding: 12px;">${cert.institution}</td>
                        <td style="padding: 12px;">${cert.course}</td>
                        <td style="padding: 12px;">${date}</td>
                        <td style="padding: 12px;"><code style="font-size: 11px; color: #666;">${cert.file_hash.substring(0, 15)}...</code></td>
                    </tr>
                `;
                tbody.innerHTML += row;
            });
        } else {
            tbody.innerHTML = "<tr><td colspan='5' style='text-align:center; padding:20px;'>Nenhum certificado encontrado.</td></tr>";
        }
    } catch (err) {
        console.error("Erro ao carregar dashboard:", err);
    }
}

// NOVA: Função de filtragem dinâmica
function filterTable() {
    const input = document.getElementById("searchInput");
    const filter = input.value.toUpperCase();
    const table = document.getElementById("certTable");
    const tr = table.getElementsByTagName("tr");

    // Percorre todas as linhas da tabela (exceto o cabeçalho)
    for (let i = 1; i < tr.length; i++) {
        let visible = false;
        const tds = tr[i].getElementsByTagName("td");
        
        // Testa se o filtro bate com o Aluno (coluna 0) ou Curso (coluna 2)
        if (tds[0] || tds[2]) {
            const nameText = tds[0].textContent || tds[0].innerText;
            const courseText = tds[2].textContent || tds[2].innerText;
            if (nameText.toUpperCase().indexOf(filter) > -1 || courseText.toUpperCase().indexOf(filter) > -1) {
                visible = true;
            }
        }
        tr[i].style.display = visible ? "" : "none";
    }
}

// Garante que carrega ao iniciar
window.onload = loadDashboard;

async function loadDashboard() {
    document.getElementById("countCerts").innerText = certs.length; // Atualiza o contador
    const tbody = document.getElementById('dashboardBody');
    try {
        const response = await fetch('/list');
        const certs = await response.json();
        tbody.innerHTML = ""; 

        certs.forEach(cert => {
            const date = new Date(cert.timestamp * 1000).toLocaleDateString();
            const row = `
                <tr>
                    <td>
                        <div style="font-weight: 600;">${cert.student_name}</div>
                        <div style="font-size: 0.75rem; color: #94a3b8;">ID: ${cert.id.substring(0,8)}</div>
                    </td>
                    <td>
                        <div>${cert.course}</div>
                        <div style="font-size: 0.8rem; color: #64748b;">${cert.institution}</div>
                    </td>
                    <td>${date}</td>
                    <td>
                        <button class="btn-verify" onclick="copyHash('${cert.file_hash}')">Copiar Hash</button>
                    </td>
                </tr>
            `;
            tbody.innerHTML += row;
        });
    } catch (err) {
        console.error("Erro ao carregar:", err);
    }
}

function copyHash(hash) {
    navigator.clipboard.writeText(hash);
    alert("Hash copiado para a área de transferência! Use-o no campo de verificação.");
}

function filterTable() {
    const input = document.getElementById("searchInput").value.toUpperCase();
    const rows = document.getElementById("dashboardBody").getElementsByTagName("tr");

    for (let row of rows) {
        const text = row.innerText.toUpperCase();
        row.style.display = text.includes(input) ? "" : "none";
    }
}