document.addEventListener("DOMContentLoaded", function () {

    const tahunSelect = document.getElementById("tahun");
    const periodeSelect = document.getElementById("periode");
    const institusiSelect = document.getElementById("institusi");
    const batchSelect = document.getElementById("batch");
    const filterForm = document.getElementById("filterForm");

    function fillSelect(selectElement, data, placeholder) {

        selectElement.innerHTML = "";

        const defaultOption = document.createElement("option");

        defaultOption.value = "";
        defaultOption.textContent = placeholder;

        selectElement.appendChild(defaultOption);


        data.forEach(function (item) {

            const option = document.createElement("option");

            option.value = item;
            option.textContent = item;

            selectElement.appendChild(option);

        });
    }

    async function loadTahun() {

        try {

            const response = await fetch("/api/filter/tahun");

            if (!response.ok) {
                throw new Error("Gagal mengambil data tahun.");
            }

            const result = await response.json();

            console.log("Data tahun:", result);

            fillSelect(
                tahunSelect,
                result.data || [],
                "Pilih Tahun"
            );

        } catch (error) {

            console.error(error);

            alert("Gagal mengambil data tahun.");

        }
    }

    tahunSelect.addEventListener("change", async function () {

        const tahun = this.value;

        console.log("Tahun dipilih:", tahun);

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


        try {

            const response = await fetch(
                `/api/filter/periode?tahun=${encodeURIComponent(tahun)}`
            );

            console.log(
                "Status API periode:",
                response.status
            );


            if (!response.ok) {
                throw new Error("Gagal mengambil periode.");
            }


            const result = await response.json();

            console.log("Data periode:", result);


            fillSelect(
                periodeSelect,
                result.data || [],
                "Pilih Periode"
            );

            periodeSelect.disabled = false;


        } catch (error) {

            console.error(error);

            alert("Gagal mengambil data periode.");

        }

    });

    periodeSelect.addEventListener("change", async function () {

        const tahun = tahunSelect.value;
        const periode = this.value;

        console.log(
            "Tahun:",
            tahun,
            "Periode:",
            periode
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


        institusiSelect.disabled = true;
        batchSelect.disabled = true;


        if (!tahun || !periode) {
            return;
        }


        try {

            const response = await fetch(
                `/api/filter/institusi?tahun=${encodeURIComponent(tahun)}&periode=${encodeURIComponent(periode)}`
            );


            console.log(
                "Status API institusi:",
                response.status
            );


            if (!response.ok) {
                throw new Error("Gagal mengambil institusi.");
            }


            const result = await response.json();

            console.log("Data institusi:", result);


            fillSelect(
                institusiSelect,
                result.data || [],
                "Pilih Asal Institusi"
            );


            institusiSelect.disabled = false;


        } catch (error) {

            console.error(error);

            alert("Gagal mengambil data institusi.");

        }

    });

    institusiSelect.addEventListener("change", async function () {

        const tahun = tahunSelect.value;
        const periode = periodeSelect.value;
        const institusi = this.value;

        console.log(
            "Tahun:",
            tahun,
            "Periode:",
            periode,
            "Institusi:",
            institusi
        );

        fillSelect(
            batchSelect,
            [],
            "Pilih Batch"
        );

        batchSelect.disabled = true;


        if (!tahun || !periode || !institusi) {
            return;
        }


        try {

            const response = await fetch(
                `/api/filter/batch?tahun=${encodeURIComponent(tahun)}&periode=${encodeURIComponent(periode)}&institusi=${encodeURIComponent(institusi)}`
            );


            console.log(
                "Status API batch:",
                response.status
            );


            if (!response.ok) {
                throw new Error("Gagal mengambil batch.");
            }


            const result = await response.json();

            console.log("Data batch:", result);


            fillSelect(
                batchSelect,
                result.data || [],
                "Pilih Batch"
            );


            batchSelect.disabled = false;


        } catch (error) {

            console.error(error);

            alert("Gagal mengambil data batch.");

        }

    });

    filterForm.addEventListener("submit", function (event) {

        event.preventDefault();


        const tahun = tahunSelect.value;
        const periode = periodeSelect.value;
        const institusi = institusiSelect.value;
        const batch = batchSelect.value;


        if (!tahun || !periode || !institusi || !batch) {

            alert(
                "Silakan pilih Tahun, Periode, Asal Institusi, dan Batch."
            );

            return;
        }


        const params = new URLSearchParams({
            tahun: tahun,
            periode: periode,
            institusi: institusi,
            batch: batch
        });


        window.location.href =
            `/analisis/result?${params.toString()}`;

    });

    const btnProsesAnalisis =
        document.getElementById("btnProsesAnalisis");

    btnProsesAnalisis.addEventListener("click", async function () {

        const konfirmasi = confirm(
            "Apakah Anda yakin ingin memproses seluruh data analisis?"
        );

        if (!konfirmasi) {
            return;
        }

        try {

            btnProsesAnalisis.disabled = true;
            btnProsesAnalisis.textContent = "Memproses...";


            const response = await fetch(
                "/analisis/process",
                {
                    method: "POST"
                }
            );


            const result = await response.json();


            if (!response.ok) {

                throw new Error(
                    result.message ||
                    "Gagal memproses analisis."
                );

            }


            alert(
                result.message ||
                "Analisis berhasil diproses."
            );


        } catch (error) {

            console.error(error);

            alert(
                error.message ||
                "Terjadi kesalahan saat memproses analisis."
            );


        } finally {

            btnProsesAnalisis.disabled = false;

            btnProsesAnalisis.textContent =
                "Proses Analisis";

        }

    });

    loadTahun();

});