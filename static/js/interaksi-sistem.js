(function () {
    "use strict";

    function getToastContainer() {
        let container = document.getElementById("toastContainer");

        if (!container) {
            container = document.createElement("div");
            container.id = "toastContainer";
            container.className = "toast-container";
            document.body.appendChild(container);
        }

        return container;
    }

    function showToast(message, type) {
        type = type || "success";

        const container = getToastContainer();

        const toast = document.createElement("div");
        toast.className = "toast toast-" + type;
        toast.textContent = message;
        container.appendChild(toast);

        requestAnimationFrame(function () {
            toast.classList.add("toast-visible");
        });

        setTimeout(function () {
            toast.classList.remove("toast-visible");
            setTimeout(function () {
                toast.remove();
            }, 300);
        }, 3500);
    }

    function getModalOverlay() {
        let overlay = document.getElementById("appModalOverlay");

        if (!overlay) {
            overlay = document.createElement("div");
            overlay.className = "app-modal-overlay";
            overlay.id = "appModalOverlay";
            overlay.hidden = true;
            overlay.innerHTML =
                '<div class="app-modal" role="dialog" aria-modal="true" aria-labelledby="appModalMessage">' +
                    '<p class="app-modal-message" id="appModalMessage"></p>' +
                    '<div class="app-modal-actions">' +
                        '<button type="button" class="btn-secondary" id="appModalCancel">Batal</button>' +
                        '<button type="button" class="btn-apply" id="appModalConfirm">Ya</button>' +
                    "</div>" +
                "</div>";
            document.body.appendChild(overlay);
        }

        return overlay;
    }

    function showConfirmModal(message) {
        return new Promise(function (resolve) {
            const overlay = getModalOverlay();
            const messageEl = document.getElementById("appModalMessage");
            const btnConfirm = document.getElementById("appModalConfirm");
            const btnCancel = document.getElementById("appModalCancel");

            messageEl.textContent = message;
            overlay.hidden = false;

            function cleanup(result) {
                overlay.hidden = true;
                btnConfirm.removeEventListener("click", onConfirm);
                btnCancel.removeEventListener("click", onCancel);
                resolve(result);
            }

            function onConfirm() {
                cleanup(true);
            }

            function onCancel() {
                cleanup(false);
            }

            btnConfirm.addEventListener("click", onConfirm);
            btnCancel.addEventListener("click", onCancel);
        });
    }

    window.showToast = showToast;
    window.showConfirmModal = showConfirmModal;

    document.addEventListener("DOMContentLoaded", function () {

        try {
            const pending = sessionStorage.getItem("pendingToast");

            if (pending) {
                sessionStorage.removeItem("pendingToast");
                showToast(pending);
            }
        } catch (e) {
            
        }

        const btnProses = document.getElementById("btnProsesAnalisis");

        if (!btnProses) {
            return;
        }

        btnProses.addEventListener("click", async function () {

            const konfirmasi = await showConfirmModal(
                "Apakah yakin ingin memproses seluruh data?"
            );

            if (!konfirmasi) {
                return;
            }

            const labelEl =
                btnProses.querySelector(".btn-proses-label") || btnProses;
            const labelAwal = labelEl.textContent;

            try {

                btnProses.disabled = true;
                btnProses.classList.add("is-loading");
                labelEl.textContent = "Memproses...";

                const response = await fetch("/analisis/process", {
                    method: "POST"
                });

                const result = await response.json();

                if (!response.ok) {
                    throw new Error(
                        result.message || "Gagal memproses analisis."
                    );
                }

                const message =
                    result.message || "Analisis berhasil diproses.";

                const lanjut = document.dispatchEvent(
                    new CustomEvent("analisis:processed", {
                        detail: { message: message },
                        cancelable: true
                    })
                );

                if (lanjut) {
                    showToast(message);
                }

            } catch (error) {

                console.error(error);

                showToast(
                    error.message ||
                    "Terjadi kesalahan saat memproses analisis.",
                    "error"
                );

            } finally {

                btnProses.disabled = false;
                btnProses.classList.remove("is-loading");
                labelEl.textContent = labelAwal;

            }

        });

    });

})();