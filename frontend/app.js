const API_BASE = 'http://localhost:8080/api';
const PROTECTED_API = `${API_BASE}/protected`;
let accessToken = localStorage.getItem('accessToken');

function updateAuthStatus() {
    const statusEl = document.getElementById('authStatus');
    if (accessToken) {
        statusEl.innerHTML = '<span class="success">Authenticated</span>';
    } else {
        statusEl.innerHTML = '<span class="error">Not authenticated</span>';
    }
}

function showError(elementId, message) {
    document.getElementById(elementId).innerHTML = `<span class="error">Error: ${message}</span>`;
}

function clearResult(elementId) {
    document.getElementById(elementId).innerHTML = '';
}

async function makeAuthRequest(url, method, body) {
    const options = {
        method,
        headers: {
            'Content-Type': 'application/json',
        },
    };
    
    if (body) {
        options.body = JSON.stringify(body);
    }
    
    const response = await fetch(url, options);
    return response;
}

async function makeRequest(url, method, body) {
    const options = {
        method,
        headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${accessToken}`
        },
    };
    
    if (body) {
        options.body = JSON.stringify(body);
    }
    
    const response = await fetch(url, options);
    if (response.status === 401) {
        // Token expired
        accessToken = null;
        localStorage.removeItem('accessToken');
        updateAuthStatus();
        throw new Error('Unauthorized');
    }
    return response;
}

// Auth functions
async function register() {
    clearResult('authStatus');
    const email = document.getElementById('email').value;
    const password = document.getElementById('password').value;
    
    try {
        const response = await makeAuthRequest(`${API_BASE}/register`, 'POST', { email, password });
        if (response.ok) {
            document.getElementById('authStatus').innerHTML = '<span class="success">Registration successful</span>';
        } else {
            const error = await response.json();
            showError('authStatus', error.error);
        }
    } catch (error) {
        showError('authStatus', error.message);
    }
}

async function login() {
    clearResult('authStatus');
    const email = document.getElementById('email').value;
    const password = document.getElementById('password').value;
    
    try {
        const response = await makeAuthRequest(`${API_BASE}/login`, 'POST', { email, password });
        if (response.ok) {
            const data = await response.json();
            accessToken = data.access_token;
            localStorage.setItem('accessToken', accessToken);
            updateAuthStatus();
        } else {
            const error = await response.json();
            showError('authStatus', error.error);
        }
    } catch (error) {
        showError('authStatus', error.message);
    }
}

function logout() {
    accessToken = null;
    localStorage.removeItem('accessToken');
    updateAuthStatus();
}

// Document functions
async function createDocument() {
    clearResult('docResult');
    const title = document.getElementById('docTitle').value;
    const sheetsCount = parseInt(document.getElementById('docSheets').value);
    const folderId = document.getElementById('docFolderId').value;
    const documentTypeId = parseInt(document.getElementById('docTypeId').value);
    
    const body = {
        title,
        sheets_count: sheetsCount,
        document_type_id: documentTypeId
    };
    
    if (folderId) {
        body.folder_id = parseInt(folderId);
    }
    
    try {
        const response = await makeRequest(`${PROTECTED_API}/documents`, 'POST', body);
        if (response.ok) {
            const doc = await response.json();
            document.getElementById('docResult').innerHTML = `
                <span class="success">Document created!</span>
                <pre>${JSON.stringify(doc, null, 2)}</pre>
            `;
        } else {
            const error = await response.json();
            showError('docResult', error.error);
        }
    } catch (error) {
        showError('docResult', error.message);
    }
}

async function getDocument() {
    clearResult('docResult');
    const docId = document.getElementById('getDocId').value;
    
    try {
        const response = await makeRequest(`${PROTECTED_API}/documents/${docId}`, 'GET');
        if (response.ok) {
            const doc = await response.json();
            document.getElementById('docResult').innerHTML = `
                <pre>${JSON.stringify(doc, null, 2)}</pre>
            `;
        } else {
            const error = await response.json();
            showError('docResult', error.error);
        }
    } catch (error) {
        showError('docResult', error.message);
    }
}

// Recommendation function
async function getRecommendation() {
    clearResult('recResult');
    const docTypeId = document.getElementById('recDocTypeId').value;
    const sheetsCount = document.getElementById('recSheetsCount').value;
    
    try {
        const params = new URLSearchParams({
            document_type_id: docTypeId,
            sheets_count: sheetsCount || '0'
        });
        
        const response = await makeRequest(`${PROTECTED_API}/folders/recommended?${params}`, 'GET');
        if (response.ok) {
            const folder = await response.json();
            if (folder) {
                document.getElementById('recResult').innerHTML = `
                    <span class="success">Recommended folder found:</span>
                    <pre>${JSON.stringify(folder, null, 2)}</pre>
                `;
            } else {
                document.getElementById('recResult').innerHTML = '<span>No suitable folder found</span>';
            }
        } else {
            const error = await response.json();
            showError('recResult', error.error);
        }
    } catch (error) {
        showError('recResult', error.message);
    }
}

// Initialize
document.addEventListener('DOMContentLoaded', () => {
    updateAuthStatus();
});
