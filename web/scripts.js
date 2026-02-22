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