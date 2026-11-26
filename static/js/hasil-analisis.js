document.addEventListener("DOMContentLoaded", function () {
    const tahunSelect = document.getElementById("tahun");
    const periodeSelect = document.getElementById("periode");
    const institusiSelect = document.getElementById("institusi");
    const batchSelect = document.getElementById("batch");
    const filterForm = document.getElementById("filterForm");
    const currentTahun = tahunSelect.dataset.selected || "";
    const currentPeriode = periodeSelect.dataset.selected || "";
    const currentInstitusi = institusiSelect.dataset.selected || "";
    const currentBatch = batchSelect.dataset.selected || "";

    function fillSelect(selectElement, data, placeholder, selectedValue = "") {
        selectElement.innerHTML = "";

        const defaultOption = document.createElement("option");
        defaultOption.value = "";
        defaultOption.textContent = placeholder;
        selectElement.appendChild(defaultOption);

        data.forEach(function (item) {
            const option = document.createElement("option");
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
            const response = await fetch("/api/filter/tahun");

            if (!response.ok) {
                throw new Error("Gagal mengambil data tahun.");
            }

            const result = await response.json();

            fillSelect(tahunSelect, result.data || [], "Pilih Tahun", selectedValue);

        } catch (error) {
            console.error(error);
            showToast("Gagal mengambil data tahun.", "error");
        }
    }

    async function loadPeriode(tahun, selectedValue = "") {
        if (!tahun) {
            fillSelect(periodeSelect, [], "Pilih Periode");
            periodeSelect.disabled = true;
            return;
        }

        try {
            const response = await fetch(
                `/api/filter/periode?tahun=${encodeURIComponent(tahun)}`
            );

            if (!response.ok) {
                throw new Error("Gagal mengambil periode.");
            }

            const result = await response.json();

            fillSelect(periodeSelect, result.data || [], "Pilih Periode", selectedValue);
            periodeSelect.disabled = false;

        } catch (error) {
            console.error(error);
            showToast("Gagal mengambil data periode.", "error");
        }
    }

    async function loadInstitusi(tahun, periode, selectedValue = "") {
        if (!tahun || !periode) {
            fillSelect(institusiSelect, [], "Pilih Asal Institusi");
            institusiSelect.disabled = true;
            return;
        }

        try {
            const response = await fetch(
                `/api/filter/institusi?tahun=${encodeURIComponent(tahun)}&periode=${encodeURIComponent(periode)}`
            );

            if (!response.ok) {
                throw new Error("Gagal mengambil institusi.");
            }

            const result = await response.json();

            fillSelect(institusiSelect, result.data || [], "Pilih Asal Institusi", selectedValue);
            institusiSelect.disabled = false;

        } catch (error) {
            console.error(error);
            showToast("Gagal mengambil data institusi.", "error");
        }
    }

    async function loadBatch(tahun, periode, institusi, selectedValue = "") {
        if (!tahun || !periode || !institusi) {
            fillSelect(batchSelect, [], "Pilih Batch");
            batchSelect.disabled = true;
            return;
        }

        try {
            const response = await fetch(
                `/api/filter/batch?tahun=${encodeURIComponent(tahun)}&periode=${encodeURIComponent(periode)}&institusi=${encodeURIComponent(institusi)}`
            );

            if (!response.ok) {
                throw new Error("Gagal mengambil batch.");
            }

            const result = await response.json();

            fillSelect(batchSelect, result.data || [], "Pilih Batch", selectedValue);
            batchSelect.disabled = false;

        } catch (error) {
            console.error(error);
            showToast("Gagal mengambil data batch.", "error");
        }
    }

    tahunSelect.addEventListener("change", async function () {
        const tahun = this.value;

        fillSelect(periodeSelect, [], "Pilih Periode");
        fillSelect(institusiSelect, [], "Pilih Asal Institusi");
        fillSelect(batchSelect, [], "Pilih Batch");

        periodeSelect.disabled = true;
        institusiSelect.disabled = true;
        batchSelect.disabled = true;

        if (!tahun) {
            return;
        }

        await loadPeriode(tahun);
    });

    periodeSelect.addEventListener("change", async function () {
        const tahun = tahunSelect.value;
        const periode = this.value;

        fillSelect(institusiSelect, [], "Pilih Asal Institusi");
        fillSelect(batchSelect, [], "Pilih Batch");

        institusiSelect.disabled = true;
        batchSelect.disabled = true;

        if (!tahun || !periode) {
            return;
        }

        await loadInstitusi(tahun, periode);
    });

    institusiSelect.addEventListener("change", async function () {
        const tahun = tahunSelect.value;
        const periode = periodeSelect.value;
        const institusi = this.value;

        fillSelect(batchSelect, [], "Pilih Batch");
        batchSelect.disabled = true;

        if (!tahun || !periode || !institusi) {
            return;
        }

        await loadBatch(tahun, periode, institusi);
    });

    filterForm.addEventListener("submit", function (event) {
        event.preventDefault();

        const tahun = tahunSelect.value;
        const periode = periodeSelect.value;
        const institusi = institusiSelect.value;
        const batch = batchSelect.value;

        if (!tahun || !periode || !institusi || !batch) {
            showToast("Silakan pilih parameter terlebih dahulu.", "error");
            return;
        }

        const params = new URLSearchParams();

        params.set("tahun", tahun);
        params.set("periode", periode);
        params.set("institusi", institusi);
        params.set("batch", batch);

        window.location.href = `/analisis/result?${params.toString()}`;
    });

    const btnApply = filterForm.querySelector(".btn-terapkan");

    function updateApplyState() {
        if (!btnApply) return;

        btnApply.disabled = !(
            tahunSelect.value &&
            periodeSelect.value &&
            institusiSelect.value &&
            batchSelect.value
        );
    }

    [tahunSelect, periodeSelect, institusiSelect, batchSelect].forEach(function (select) {
        select.addEventListener("change", updateApplyState);
    });

    const btnReset = document.getElementById("btnReset");

    if (btnReset) {
        btnReset.addEventListener("click", function () {
            if (document.querySelector(".analysis-result")) {
                window.location.href = "/analisis";
                return;
            }

            tahunSelect.value = "";

            fillSelect(periodeSelect, [], "Pilih Periode");
            fillSelect(institusiSelect, [], "Pilih Asal Institusi");
            fillSelect(batchSelect, [], "Pilih Batch");

            periodeSelect.disabled = true;
            institusiSelect.disabled = true;
            batchSelect.disabled = true;

            updateApplyState();
        });
    }

    async function initializeFilter() {
        await loadTahun(currentTahun);

        if (!currentTahun) {
            return;
        }

        await loadPeriode(currentTahun, currentPeriode);

        if (!currentPeriode) {
            return;
        }

        await loadInstitusi(currentTahun, currentPeriode, currentInstitusi);

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

    initializeFilter().then(updateApplyState);
    document.addEventListener("analisis:processed", function (event) {
        event.preventDefault(); 

        try {
            sessionStorage.setItem("pendingToast", event.detail.message);
        } catch (e) {
            
        }

        window.location.reload();
    });

    const tabs = Array.from(document.querySelectorAll(".tab-btn"));

    function activateTab(tab) {
        tabs.forEach(function (t) {
            const active = t === tab;

            t.classList.toggle("active", active);
            t.setAttribute("aria-selected", active);
            t.tabIndex = active ? 0 : -1;

            const panel = document.getElementById(t.getAttribute("aria-controls"));
            panel.hidden = !active;
            panel.classList.toggle("active", active);
        });
    }

    tabs.forEach(function (tab, i) {
        tab.addEventListener("click", function () {
            activateTab(tab);
        });

        tab.addEventListener("keydown", function (e) {
            if (e.key !== "ArrowRight" && e.key !== "ArrowLeft") return;

            const step = e.key === "ArrowRight" ? 1 : -1;
            const next = tabs[(i + step + tabs.length) % tabs.length];

            activateTab(next);
            next.focus();
        });
    });

    const pagers = new Map();

    function initPagination(panel) {
        const tbody = panel.querySelector(".result-table tbody");
        const footer = panel.querySelector(".pagination");
        if (!tbody || !footer) return;

        const rows = Array.from(tbody.querySelectorAll("tr")).filter(function (row) {
            return !row.querySelector(".empty-row");
        });

        const info = footer.querySelector(".pagination-info");
        const sizeSelect = footer.querySelector(".pagination-select");
        const prevBtn = footer.querySelector('[data-dir="prev"]');
        const nextBtn = footer.querySelector('[data-dir="next"]');

        let page = 1;

        function render() {
            const total = rows.length;
            const size = Number(sizeSelect.value);
            const pages = Math.max(1, Math.ceil(total / size));

            page = Math.min(page, pages);

            const start = (page - 1) * size;
            const end = Math.min(start + size, total);

            rows.forEach(function (row, i) {
                row.hidden = !(i >= start && i < end);
            });

            info.textContent =
                total === 0
                    ? "Menampilkan 0 dari 0 peserta"
                    : `Menampilkan ${start + 1}-${end} dari ${total} peserta`;

            prevBtn.disabled = page <= 1;
            nextBtn.disabled = page >= pages;
        }

        prevBtn.addEventListener("click", function () {
            page--;
            render();
        });

        nextBtn.addEventListener("click", function () {
            page++;
            render();
        });

        sizeSelect.addEventListener("change", function () {
            page = 1;
            render();
        });

        pagers.set(panel, {
            showAll: function () {
                rows.forEach(function (row) {
                    row.hidden = false;
                });
            },
            render: render
        });

        render();
    }

    document.querySelectorAll(".tab-panel").forEach(initPagination);

    async function captureElement(element) {
        const wasHidden =
            element.classList.contains("tab-panel") && element.hidden;

        const toHide = element.querySelectorAll(
            ".analysis-card-header, .tab-meta, .pagination"
        );

        const pager = pagers.get(element);

        if (wasHidden) element.hidden = false;
        if (pager) pager.showAll();
        toHide.forEach(function (el) {
            el.style.display = "none";
        });

        try {
            return await html2canvas(element, {
                backgroundColor: "#ffffff",
                scale: 2
            });
        } finally {
            toHide.forEach(function (el) {
                el.style.display = "";
            });
            if (pager) pager.render();
            if (wasHidden) element.hidden = true;
        }
    }

    function getCaptureName(element, index) {
        let title = "";

        if (element.classList.contains("tab-panel")) {
            const tab = document.getElementById(
                element.getAttribute("aria-labelledby")
            );
            title = tab && tab.firstChild ? tab.firstChild.textContent.trim() : "";
        } else {
            const h3 = element.querySelector("h3");
            title = h3 ? h3.textContent.trim() : "";
        }

        return title
            ? title.toLowerCase().replace(/\s+/g, "_")
            : `analisis_${index + 1}`;
    }

    const copyButtons = document.querySelectorAll(".copy-button");

    copyButtons.forEach(function (button) {
        button.addEventListener("click", async function () {
            const target =
                this.dataset.copy === "tabs"
                    ? document.querySelector(".tab-panel:not([hidden])")
                    : this.closest(".analysis-card");

            if (!target) return;

            try {
                const canvas = await captureElement(target);

                canvas.toBlob(async function (blob) {
                    if (!blob) {
                        showToast("Gagal membuat gambar.", "error");
                        return;
                    }

                    try {
                        const item = new ClipboardItem({ "image/png": blob });
                        await navigator.clipboard.write([item]);
                        showToast("Berhasil menyalin gambar", "success");
                    } catch (err) {
                        console.error("Gagal menyalin gambar:", err);
                        showToast("Gagal menyalin gambar!", "error");
                    }
                }, "image/png");

            } catch (err) {
                console.error("Gagal menyalin gambar:", err);
                showToast("Gagal menyalin gambar!", "error");
            }
        });
    });

    const downloadToggle = document.getElementById("downloadToggle");
    const downloadList = document.getElementById("downloadList");
    const downloadLabel = downloadToggle
        ? downloadToggle.querySelector(".download-label")
        : null;

    function getDownloadTargets() {
        return {
            zona: document.querySelector(".analysis-card:not(.tabs-card)"),
            kehadiran: document.getElementById("panel-attendance"),
            to: document.getElementById("panel-to")
        };
    }

    function setDownloadMenu(open) {
        downloadList.hidden = !open;
        downloadToggle.setAttribute("aria-expanded", open);
    }

    function triggerDownload(blob, fileName) {
        const link = document.createElement("a");
        link.href = URL.createObjectURL(blob);
        link.download = fileName;
        document.body.appendChild(link);
        link.click();
        document.body.removeChild(link);

        URL.revokeObjectURL(link.href);
    }

    function getZipName() {
        const elTahun = document.getElementById("tahun");
        const elPeriode = document.getElementById("periode");
        const elInstitusi = document.getElementById("institusi");
        const elBatch = document.getElementById("batch");

        const tahun = elTahun?.getAttribute("data-selected") || elTahun?.value || "";
        const periode = elPeriode?.getAttribute("data-selected") || elPeriode?.value || "";
        const asalInstitusi = elInstitusi?.getAttribute("data-selected") || elInstitusi?.value || "";
        const batch = elBatch?.getAttribute("data-selected") || elBatch?.value || "";

        const rawZipName = `Hasil Analisis Tahun ${tahun} Periode ${periode} ${asalInstitusi} Batch ${batch}`;

        return rawZipName.replace(/[/\\?%*:|"<>]/g, "").trim() + ".zip";
    }

    async function downloadZip(targets) {
        const zip = new JSZip();
        const imgFolder = zip.folder("analisis_gambar");

        const elements = [targets.zona, targets.kehadiran, targets.to];

        for (let index = 0; index < elements.length; index++) {
            const element = elements[index];
            if (!element) continue;

            const canvas = await captureElement(element);

            const imageData = canvas
                .toDataURL("image/png")
                .replace(/^data:image\/(png|jpg);base64,/, "");

            imgFolder.file(
                `${getCaptureName(element, index)}.png`,
                imageData,
                { base64: true }
            );
        }

        const content = await zip.generateAsync({ type: "blob" });

        triggerDownload(content, getZipName());
    }

    async function downloadSingle(element, index) {
        if (!element) return;

        const canvas = await captureElement(element);

        const blob = await new Promise(function (resolve) {
            canvas.toBlob(resolve, "image/png");
        });

        if (!blob) {
            throw new Error("Gagal membuat gambar.");
        }

        triggerDownload(blob, `${getCaptureName(element, index)}.png`);
    }

    if (downloadToggle && downloadList) {

        downloadToggle.addEventListener("click", function (event) {
            event.stopPropagation();
            setDownloadMenu(downloadList.hidden);
        });

        document.addEventListener("click", function (event) {
            if (!downloadList.hidden && !downloadList.contains(event.target)) {
                setDownloadMenu(false);
            }
        });

        document.addEventListener("keydown", function (event) {
            if (event.key === "Escape" && !downloadList.hidden) {
                setDownloadMenu(false);
                downloadToggle.focus();
            }
        });

        downloadList.addEventListener("click", async function (event) {
            const item = event.target.closest("[data-download]");
            if (!item) return;

            setDownloadMenu(false);

            const targets = getDownloadTargets();
            const kind = item.dataset.download;
            const order = { zona: 0, kehadiran: 1, to: 2 };

            downloadToggle.disabled = true;
            const originalText = downloadLabel.textContent;
            downloadLabel.textContent = "Memproses...";

            try {
                if (kind === "all") {
                    await downloadZip(targets);
                } else {
                    await downloadSingle(targets[kind], order[kind]);
                }
            } catch (err) {
                console.error("Gagal mengunduh:", err);
                alert("Terjadi kesalahan saat mengunduh.");
            } finally {
                downloadToggle.disabled = false;
                downloadLabel.textContent = originalText;
            }
        });
    }

    const zonaCanvas = document.getElementById("zonaChart");

    if (zonaCanvas) {
        const zona = {
            merah: Number(zonaCanvas.dataset.merah || 0),
            kuning: Number(zonaCanvas.dataset.kuning || 0),
            hijau: Number(zonaCanvas.dataset.hijau || 0)
        };

        const total = zona.merah + zona.kuning + zona.hijau;

        // legend: "Zona merah · 5 (3%)"
        Object.keys(zona).forEach(function (key) {
            const el = document.querySelector(`[data-zona-detail="${key}"]`);
            if (!el) return;

            const pct = total > 0 ? Math.round((zona[key] / total) * 100) : 0;
            el.textContent = `· ${zona[key]} (${pct}%)`;
        });

        const zonaValueLabel = {
            id: "zonaValueLabel",
            afterDatasetsDraw: function (chart) {
                if (total === 0) return;

                const ctx = chart.ctx;
                const meta = chart.getDatasetMeta(0);
                const values = chart.data.datasets[0].data;

                ctx.save();
                ctx.fillStyle = "#ffffff";
                ctx.font = "700 13px 'Plus Jakarta Sans', sans-serif";
                ctx.textAlign = "center";
                ctx.textBaseline = "middle";

                meta.data.forEach(function (arc, i) {
                    if (values[i] / total < 0.1) return;

                    const pos = arc.tooltipPosition();
                    ctx.fillText(values[i], pos.x, pos.y);
                });

                ctx.restore();
            }
        };

        if (typeof Chart !== "undefined") {
            new Chart(zonaCanvas, {
                type: "doughnut",

                data: {
                    labels:
                        total > 0
                            ? ["Zona merah", "Zona kuning", "Zona hijau"]
                            : ["Belum ada data"],
                    datasets: [
                        {
                            data:
                                total > 0
                                    ? [zona.merah, zona.kuning, zona.hijau]
                                    : [1],
                            backgroundColor:
                                total > 0
                                    ? ["#E24C4B", "#FFC700", "#2FB380"]
                                    : ["#EEF2F5"],
                            borderColor: "#ffffff",
                            borderWidth: 3,
                            hoverOffset: 0
                        }
                    ]
                },

                options: {
                    responsive: true,
                    maintainAspectRatio: false,
                    cutout: "70%",
                    plugins: {
                        legend: { display: false },
                        tooltip: { enabled: total > 0 }
                    }
                },

                plugins: [zonaValueLabel]
            });
        }
    }

});