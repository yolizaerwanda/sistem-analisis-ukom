document.addEventListener("DOMContentLoaded", function () {

    const tahunSelect = document.getElementById("tahun");
    const periodeSelect = document.getElementById("periode");
    const institusiSelect = document.getElementById("institusi");
    const batchSelect = document.getElementById("batch");
    const filterForm = document.getElementById("filterForm");

    const currentTahun =
        tahunSelect.dataset.selected || "";

    const currentPeriode =
        periodeSelect.dataset.selected || "";

    const currentInstitusi =
        institusiSelect.dataset.selected || "";

    const currentBatch =
        batchSelect.dataset.selected || "";

    function fillSelect(
        selectElement,
        data,
        placeholder,
        selectedValue = ""
    ) {

        selectElement.innerHTML = "";

        const defaultOption =
            document.createElement("option");

        defaultOption.value = "";
        defaultOption.textContent = placeholder;

        selectElement.appendChild(defaultOption);


        data.forEach(function (item) {

            const option =
                document.createElement("option");

            option.value = item;
            option.textContent = item;

            if (item === selectedValue) {
                option.selected = true;
            }

            selectElement.appendChild(option);

        });

    }

    async function loadTahun(selectedValue = "") {

        try {

            const response =
                await fetch("/api/filter/tahun");

            if (!response.ok) {
                throw new Error(
                    "Gagal mengambil data tahun."
                );
            }

            const result =
                await response.json();

            fillSelect(
                tahunSelect,
                result.data || [],
                "Pilih Tahun",
                selectedValue
            );

        } catch (error) {

            console.error(error);

            showToast(
                "Gagal mengambil data tahun.",
                "error"
            );

        }

    }

    async function loadPeriode(
        tahun,
        selectedValue = ""
    ) {

        if (!tahun) {

            fillSelect(
                periodeSelect,
                [],
                "Pilih Periode"
            );

            periodeSelect.disabled = true;

            return;

        }

        try {

            const response =
                await fetch(
                    `/api/filter/periode?tahun=${encodeURIComponent(tahun)}`
                );

            if (!response.ok) {
                throw new Error(
                    "Gagal mengambil periode."
                );
            }

            const result =
                await response.json();

            fillSelect(
                periodeSelect,
                result.data || [],
                "Pilih Periode",
                selectedValue
            );

            periodeSelect.disabled = false;

        } catch (error) {

            console.error(error);

            showToast(
                "Gagal mengambil data periode.",
                "error"
            );

        }

    }

    async function loadInstitusi(
        tahun,
        periode,
        selectedValue = ""
    ) {

        if (!tahun || !periode) {

            fillSelect(
                institusiSelect,
                [],
                "Pilih Asal Institusi"
            );

            institusiSelect.disabled = true;

            return;

        }

        try {

            const response =
                await fetch(
                    `/api/filter/institusi?tahun=${encodeURIComponent(tahun)}&periode=${encodeURIComponent(periode)}`
                );

            if (!response.ok) {
                throw new Error(
                    "Gagal mengambil institusi."
                );
            }

            const result =
                await response.json();

            fillSelect(
                institusiSelect,
                result.data || [],
                "Pilih Asal Institusi",
                selectedValue
            );

            institusiSelect.disabled = false;

        } catch (error) {

            console.error(error);

            showToast(
                "Gagal mengambil data institusi.",
                "error"
            );

        }

    }

    async function loadBatch(
        tahun,
        periode,
        institusi,
        selectedValue = ""
    ) {

        if (
            !tahun ||
            !periode ||
            !institusi
        ) {

            fillSelect(
                batchSelect,
                [],
                "Pilih Batch"
            );

            batchSelect.disabled = true;

            return;

        }

        try {

            const response =
                await fetch(
                    `/api/filter/batch?tahun=${encodeURIComponent(tahun)}&periode=${encodeURIComponent(periode)}&institusi=${encodeURIComponent(institusi)}`
                );

            if (!response.ok) {
                throw new Error(
                    "Gagal mengambil batch."
                );
            }

            const result =
                await response.json();

            fillSelect(
                batchSelect,
                result.data || [],
                "Pilih Batch",
                selectedValue
            );

            batchSelect.disabled = false;

        } catch (error) {

            console.error(error);

            showToast(
                "Gagal mengambil data batch.",
                "error"
            );

        }

    }

    tahunSelect.addEventListener(
        "change",
        async function () {

            const tahun = this.value;

            fillSelect(
                periodeSelect,
                [],
                "Pilih Periode"
            );

            fillSelect(
                institusiSelect,
                [],
                "Pilih Asal Institusi"
            );

            fillSelect(
                batchSelect,
                [],
                "Pilih Batch"
            );

            periodeSelect.disabled = true;
            institusiSelect.disabled = true;
            batchSelect.disabled = true;

            if (!tahun) {
                return;
            }

            await loadPeriode(tahun);

        }
    );

    periodeSelect.addEventListener(
        "change",
        async function () {

            const tahun =
                tahunSelect.value;

            const periode =
                this.value;

            fillSelect(
                institusiSelect,
                [],
                "Pilih Asal Institusi"
            );

            fillSelect(
                batchSelect,
                [],
                "Pilih Batch"
            );

            institusiSelect.disabled = true;
            batchSelect.disabled = true;

            if (!tahun || !periode) {
                return;
            }

            await loadInstitusi(
                tahun,
                periode
            );

        }
    );

    institusiSelect.addEventListener(
        "change",
        async function () {

            const tahun =
                tahunSelect.value;

            const periode =
                periodeSelect.value;

            const institusi =
                this.value;

            fillSelect(
                batchSelect,
                [],
                "Pilih Batch"
            );

            batchSelect.disabled = true;

            if (
                !tahun ||
                !periode ||
                !institusi
            ) {
                return;
            }

            await loadBatch(
                tahun,
                periode,
                institusi
            );

        }
    );

    filterForm.addEventListener(
        "submit",
        function (event) {

            event.preventDefault();

            const tahun =
                tahunSelect.value;

            const periode =
                periodeSelect.value;

            const institusi =
                institusiSelect.value;

            const batch =
                batchSelect.value;

            if (
                !tahun ||
                !periode ||
                !institusi ||
                !batch
            ) {

                showToast(
                    "Silakan pilih parameter terlebih dahulu.",
                    "error"
                );

                return;

            }

            const params =
                new URLSearchParams();

            params.set("tahun", tahun);
            params.set("periode", periode);
            params.set("institusi", institusi);
            params.set("batch", batch);

            window.location.href =
                `/analisis/result?${params.toString()}`;

        }
    );

    async function initializeFilter() {

        await loadTahun(
            currentTahun
        );

        if (!currentTahun) {
            return;
        }

        await loadPeriode(
            currentTahun,
            currentPeriode
        );

        if (!currentPeriode) {
            return;
        }

        await loadInstitusi(
            currentTahun,
            currentPeriode,
            currentInstitusi
        );

        if (!currentInstitusi) {
            return;
        }

        await loadBatch(
            currentTahun,
            currentPeriode,
            currentInstitusi,
            currentBatch
        );

    }


    initializeFilter();

    const copyButtons = document.querySelectorAll(".copy-button");

    copyButtons.forEach(function (button) {
        button.addEventListener("click", async function () {
            const cardElement = this.closest(".analysis-card");
            if (!cardElement) return;

            const header = cardElement.querySelector(".analysis-card-header");
            const footer = cardElement.querySelector(".analysis-card-footer");

            try {
                if (header) header.style.display = "none";
                if (footer) footer.style.display = "none";

                const canvas = await html2canvas(cardElement, {
                    backgroundColor: "#ffffff",
                    scale: 2
                });

                if (header) header.style.display = "";
                if (footer) footer.style.display = "";

                canvas.toBlob(async function (blob) {
                    if (!blob) {
                        alert("Gagal membuat gambar.");
                        return;
                    }

                    const item = new ClipboardItem({ "image/png": blob });
                    await navigator.clipboard.write([item]);

                    showToast("Berhasil menyalin gambar", "success");
                }, "image/png");

            } catch (err) {
                if (header) header.style.display = "";
                if (footer) footer.style.display = "";

                console.error("Gagal menyalin gambar:", err);
                showToast("Gagal menyalin gambar!", "error");
            }
        });
    });

    function showToast(message, type = "success") {
        let container = document.getElementById("toastContainer");

        if (!container) {
            container = document.createElement("div");
            container.id = "toastContainer";
            container.className = "toast-container";
            document.body.appendChild(container);
        }

        const toast = document.createElement("div");
        toast.className = `toast toast-${type}`;
        toast.textContent = message;
        container.appendChild(toast);

        requestAnimationFrame(() => {
            toast.classList.add("toast-visible");
        });

        setTimeout(() => {
            toast.classList.remove("toast-visible");
            setTimeout(() => toast.remove(), 300);
        }, 3500);
    }

    const downloadButton = document.getElementById("downloadAll");

    if (downloadButton) {
        downloadButton.addEventListener("click", async function () {
            const cards = document.querySelectorAll(".analysis-card");

            if (cards.length === 0) {
                alert("Tidak ada hasil analisis untuk diunduh.");
                return;
            }

            downloadButton.disabled = true;
            const originalText = downloadButton.textContent;
            downloadButton.textContent = "Memproses ZIP...";

            try {
                const zip = new JSZip();
                const imgFolder = zip.folder("analisis_gambar");

                for (let index = 0; index < cards.length; index++) {
                    const cardElement = cards[index];

                    const header = cardElement.querySelector(".analysis-card-header");
                    const footer = cardElement.querySelector(".analysis-card-footer");
                    const titleElement = cardElement.querySelector("h3");

                    let fileName = titleElement 
                        ? titleElement.textContent.trim().toLowerCase().replace(/\s+/g, "_") 
                        : `analisis_${index + 1}`;

                    if (header) header.style.display = "none";
                    if (footer) footer.style.display = "none";

                    const canvas = await html2canvas(cardElement, {
                        backgroundColor: "#ffffff",
                        scale: 2
                    });

                    if (header) header.style.display = "";
                    if (footer) footer.style.display = "";

                    const imageData = canvas.toDataURL("image/png").replace(/^data:image\/(png|jpg);base64,/, "");

                    imgFolder.file(`${fileName}.png`, imageData, { base64: true });
                }

                const elTahun = document.getElementById("tahun");
                const elPeriode = document.getElementById("periode");
                const elInstitusi = document.getElementById("institusi");
                const elBatch = document.getElementById("batch");

                const tahun = elTahun?.getAttribute("data-selected") || elTahun?.value || "";
                const periode = elPeriode?.getAttribute("data-selected") || elPeriode?.value || "";
                const asalInstitusi = elInstitusi?.getAttribute("data-selected") || elInstitusi?.value || "";
                const batch = elBatch?.getAttribute("data-selected") || elBatch?.value || "";

                const rawZipName = `Hasil Analisis Tahun ${tahun} Periode ${periode} ${asalInstitusi} Batch ${batch}`;
                
                const zipFileName = rawZipName.replace(/[/\\?%*:|"<>]/g, "").trim() + ".zip";

                const content = await zip.generateAsync({ type: "blob" });
                
                const link = document.createElement("a");
                link.href = URL.createObjectURL(content);
                link.download = zipFileName;
                document.body.appendChild(link);
                link.click();
                document.body.removeChild(link);

                URL.revokeObjectURL(link.href);

            } catch (err) {
                console.error("Gagal mengunduh ZIP:", err);
                alert("Terjadi kesalahan saat membuat file ZIP.");
            } finally {
                downloadButton.disabled = false;
                downloadButton.textContent = originalText;
            }
        });
    }

    const zonaCanvas =
        document.getElementById(
            "zonaChart"
        );


    if (zonaCanvas) {

        const merah =
            Number(
                zonaCanvas.dataset.merah || 0
            );

        const kuning =
            Number(
                zonaCanvas.dataset.kuning || 0
            );

        const hijau =
            Number(
                zonaCanvas.dataset.hijau || 0
            );


        if (typeof Chart !== "undefined") {

            new Chart(
                zonaCanvas,
                {
                    type: "pie",

                    data: {

                        datasets: [
                            {
                                data: [
                                    merah,
                                    kuning,
                                    hijau
                                ],
                                backgroundColor: [
                                    "#ff4d4f", 
                                    "#ffec3d", 
                                    "#52c41a"  
                                ]
                            }
                        ]

                    },

                    options: {

                        responsive: true,

                        maintainAspectRatio: false

                    }

                }
            );

        }

    }

});