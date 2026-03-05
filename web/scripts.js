const API_URL = 'http://localhost:8080';

document.addEventListener('DOMContentLoaded', () => {
    document.getElementById('registerForm').addEventListener('submit', registerCertificate);
    loadDashboard();
});

async function registerCertificate(e) {
    e.preventDefault();
    const resultDiv = document.getElementById('registerResult');
    resultDiv.innerHTML = "<strong>Minerando bloco...</strong> Isso pode levar alguns segundos devido ao PoW.";

    const formData = new FormData(e.target);

    try {
        const response = await fetch(`${API_URL}/register`, { 
            method: 'POST', 
            body: formData 
        });
        
        const data = await response.json();
        
        if (response.ok) {
            resultDiv.innerHTML = `
                <strong style="color: #10b981;">Sucesso!</strong><br>
                ID: ${data.id}<br>
                Hash: <small>${data.hash}</small>
            `;
            e.target.reset(); // Limpa o formulário
            loadDashboard();  // Atualiza a tabela automaticamente
        } else {
            const errorText = await response.text();
            resultDiv.innerText = "Erro: " + (data.error || "Falha no registro");
        }
    } catch (err) {
        resultDiv.innerText = "Erro ao conectar com o servidor. Certifique-se que o backend esta rodando em :8080";
        console.error("erro real no fetch:", err);
    }
}

async function verifyCertificate() {
    const hash = document.getElementById('hashInput').value.trim();
    const resultDiv = document.getElementById('verifyResult');

    if (!hash) {
        resultDiv.innerText = "Por favor, insira um hash.";
        return;
    }

    try {
        const response = await fetch(`${API_URL}/verify?hash=${encodeURIComponent(hash)}`);
        const data = await response.json();
        
        if (response.ok && !data.error) {
            // Ajustado para os nomes das propriedades JSON que o Go envia
            resultDiv.innerHTML = `
                <div style="border-left: 4px solid #10b981; padding-left: 10px;">
                    <strong style="color: #10b981;">Certificado Autêntico!</strong><br>
                    <strong>Aluno:</strong> ${data.student_name}<br>
                    <strong>Instituição:</strong> ${data.institution}<br>
                    <strong>Curso:</strong> ${data.course}<br>
                    <strong>Data:</strong> ${new Date(data.timestamp * 1000).toLocaleString()}
                </div>
            `;
        } else {
            resultDiv.innerHTML = `<strong style="color: #ef4444;">Certificado não encontrado ou inválido na Blockchain.</strong>`;
        }
    } catch (err) {
        resultDiv.innerText = "Erro ao consultar servidor.";
        console.error(err);
    }
}

async function loadDashboard() {
    const tbody = document.getElementById('dashboardBody');
    const counter = document.getElementById('countCerts');

    try {
        const response = await fetch(`${API_URL}/list`);
        const certs = await response.json();

        tbody.innerHTML = "";
        if (Array.isArray(certs) && certs.length > 0) {
            certs.forEach(cert => {
                const date = cert.timestamp ? new Date(cert.timestamp * 1000).toLocaleDateString() : 'N/A';
                const row = `
                    <tr>
                        <td><strong>${cert.student_name}</strong></td>
                        <td>${cert.course} <br> <small>${cert.institution}</small></td>
                        <td>${date}</td>
                        <td><code title="${cert.file_hash}">${cert.file_hash.substring(0, 10)}...</code></td>
                        <td>
                            <button class="btn btn-primary" style="padding: 5px 10px; font-size: 0.8rem;" 
                                onclick="copyHash('${cert.file_hash}')">Copiar Hash</button>
                        </td>
                    </tr>
                `;
                tbody.insertAdjacentHTML('beforeend', row);
            });
            if (counter) counter.innerText = `(${certs.length})`;
        } else {
            tbody.innerHTML = "<tr><td colspan='5' style='text-align:center;padding:20px;'>Nenhum certificado registrado na rede ainda.</td></tr>";
            if (counter) counter.innerText = `(0)`;
        }
    } catch (err) {
        console.error("Erro ao carregar dashboard:", err);
        tbody.innerHTML = "<tr><td colspan='5' style='text-align:center;color:red;'>Erro ao carregar dados do servidor.</td></tr>";
    }
}

function filterTable() {
    const input = document.getElementById("searchInput");
    if (!input) return;
    const filter = input.value.toUpperCase();
    const rows = document.getElementById("dashboardBody").getElementsByTagName("tr");

    for (let row of rows) {
        row.style.display = row.innerText.toUpperCase().includes(filter) ? "" : "none";
    }
}

function copyHash(hash) {
    if (!hash) return;
    navigator.clipboard.writeText(hash).then(() => {
        alert("Hash copiado com sucesso!");
    });
}