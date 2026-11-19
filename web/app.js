"use strict";
const $ = (id) => document.getElementById(id);
const src = $("src"), frame = $("frame"), printFrame = $("print-frame");
const DRAFT_KEY = "wsp-draft";
let mode = "interactive", result = null, timer = 0;

const store = {
  get() { try { return localStorage.getItem(DRAFT_KEY); } catch { return null; } },
  set(v) { try { localStorage.setItem(DRAFT_KEY, v); } catch {} },
};

function slug(title) {
  return (title || "worksheet").normalize("NFKD").replace(/[^\w\s-]/g, "").trim().replace(/\s+/g, "-").toLowerCase() || "worksheet";
}

function download(name, text, type = "text/html") {
  const a = document.createElement("a");
  a.href = URL.createObjectURL(new Blob([text], { type: type + ";charset=utf-8" }));
  a.download = name;
  a.click();
  setTimeout(() => URL.revokeObjectURL(a.href), 1000);
}

function showWarnings(list, questions) {
  const box = $("warnings");
  box.replaceChildren();
  if (list.length) {
    const ul = document.createElement("ul");
    for (const w of list) ul.appendChild(Object.assign(document.createElement("li"), { textContent: w }));
    box.appendChild(ul);
  } else {
    box.appendChild(Object.assign(document.createElement("p"), { textContent: "No problems found." }));
  }
  $("count").textContent = `· ${questions} question${questions === 1 ? "" : "s"}`;
}

function update() {
  result = wsp.render(src.value);
  showWarnings(result.warnings, result.questions);
  frame.srcdoc = result[mode];
  store.set(src.value);
}

function setMode(m) {
  mode = m;
  history.replaceState(null, "", m === "printable" ? "#printable" : location.pathname);
  for (const t of ["interactive", "printable"]) $("tab-" + t).setAttribute("aria-selected", String(t === m));
  if (result) frame.srcdoc = result[m];
}

async function loadExample() {
  const res = await fetch("example.wsp.txt");
  src.value = await res.text();
  update();
}

$("tab-interactive").onclick = () => setMode("interactive");
$("tab-printable").onclick = () => setMode("printable");
src.addEventListener("input", () => { clearTimeout(timer); timer = setTimeout(update, 250); });
$("example").onclick = () => {
  if (src.value.trim() && !confirm("Replace the current text with the example?")) return;
  loadExample();
};
$("open").onchange = async (e) => {
  const file = e.target.files[0];
  if (!file) return;
  src.value = await file.text();
  e.target.value = "";
  update();
};
$("save-src").onclick = () => download(slug(result?.title) + ".wsp.txt", src.value, "text/plain");
$("download").onclick = () => result && download(`${slug(result.title)}_${mode}.html`, result[mode]);
$("print").onclick = () => {
  if (!result) return;
  // The printable sheet has no scripts, so it can live in an unsandboxed frame we can call print() on.
  printFrame.onload = () => printFrame.contentWindow.print();
  printFrame.srcdoc = result.printable;
};
$("help-btn").onclick = () => $("help").showModal();

(async () => {
  try {
    const go = new Go();
    const { instance } = await WebAssembly.instantiateStreaming(fetch("main.wasm"), go.importObject);
    go.run(instance);
  } catch (err) {
    $("loading").textContent = "Could not load the parser: " + err.message;
    return;
  }
  $("loading").hidden = true;
  if (location.hash === "#printable") setMode("printable");
  const draft = store.get();
  if (draft) { src.value = draft; update(); } else { await loadExample(); }
})();
