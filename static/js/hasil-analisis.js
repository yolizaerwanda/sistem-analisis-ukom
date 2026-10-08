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

            alert(
                "Gagal mengambil data tahun."
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

            alert(
                "Gagal mengambil data periode."
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

            alert(
                "Gagal mengambil data institusi."
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

            alert(
                "Gagal mengambil data batch."
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

                alert(
                    "Silakan pilih parameter terlebih dahulu.",
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

        button.addEventListener(
            "click",
            function () {

                alert(
                    "Fitur Copy akan dibuat pada tahap berikutnya."
                );

            }
        );

    });

    const downloadButton = document.getElementById("downloadAll");

    if (downloadButton) {

        downloadButton.addEventListener(
            "click",
            function () {

                alert(
                    "Fitur download akan dibuat pada tahap berikutnya."
                );

            }
        );

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