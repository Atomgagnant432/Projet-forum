function valider() {
    const email = document.getElementById('email').value.trim();
    const pwd = document.getElementById('pwd').value;
    const erreur = document.getElementById('error-msg');

    if (!email || !pwd) {
        erreur.textContent = 'Veuillez remplir tous les champs.';
        return;
    }

    alert('Connecté !');
}