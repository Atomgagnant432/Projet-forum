function valider() {
    const pseudo = document.getElementById('pseudo').value.trim();
    const email = document.getElementById('email').value.trim();
    const pwd = document.getElementById('pwd').value;
    const pwdConfirm = document.getElementById('pwd-confirm').value;
    const erreur = document.getElementById('error-msg');

    if (!pseudo || !email || !pwd || !pwdConfirm) {
        erreur.textContent = 'Veuillez remplir tous les champs.';
        return;
    }

    if (pwd !== pwdConfirm) {
        erreur.textContent = 'Les mots de passe ne correspondent pas.';
        return;
    }

    erreur.textContent = '';
    
    alert('Compte créé !');
} 