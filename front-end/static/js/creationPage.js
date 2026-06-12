const imageInput = document.getElementById("image");
const previewContainer = document.getElementById("image-preview-container");
const imagePreview = document.getElementById("image-preview");
const removeImageButton = document.getElementById("remove-image-btn");

let previewURL = null;

function removePreview() {
    if (previewURL !== null) {
        URL.revokeObjectURL(previewURL);
        previewURL = null;
    }

    imageInput.value = "";
    imagePreview.removeAttribute("src");
    previewContainer.hidden = true;
}

imageInput.addEventListener("change", () => {
    const file = imageInput.files[0];

    if (!file) {
        removePreview();
        return;
    }

    if (!file.type.startsWith("image/")) {
        alert("Le fichier sélectionné doit être une image.");
        removePreview();
        return;
    }

    const maximumSize = 20 * 1024 * 1024;

    if (file.size > maximumSize) {
        alert("L’image ne doit pas dépasser 20 Mo.");
        removePreview();
        return;
    }

    if (previewURL !== null) {
        URL.revokeObjectURL(previewURL);
    }

    previewURL = URL.createObjectURL(file);
    imagePreview.src = previewURL;
    previewContainer.hidden = false;
});

removeImageButton.addEventListener("click", removePreview);