// -----------------------------------------------------------------------------
// Init Bhojpur Web service worker
// -----------------------------------------------------------------------------
var bhojpurOnUpdate = function () { };

if ("serviceWorker" in navigator) {
  navigator.serviceWorker
    .register("{{.WorkerJS}}")
    .then(reg => {
      console.log("registering Bhojpur Web application service worker");

      reg.onupdatefound = function () {
        const installingWorker = reg.installing;
        installingWorker.onstatechange = function () {
          if (installingWorker.state == "installed") {
            if (navigator.serviceWorker.controller) {
              bhojpurOnUpdate();
            }
          }
        };
      }
    })
    .catch(err => {
      console.error("offline Bhojpur Web service worker registration failed", err);
    });
}

// -----------------------------------------------------------------------------
// Env
// -----------------------------------------------------------------------------
const bhojpurEnv = {{.Env }};

function bhojpurGetenv(k) {
  return bhojpurEnv[k];
}

// -----------------------------------------------------------------------------
// Bhojpur Web application install
// -----------------------------------------------------------------------------
let deferredPrompt = null;
var bhojpurOnAppInstallChange = function () { };

window.addEventListener("beforeinstallprompt", e => {
  e.preventDefault();
  deferredPrompt = e;
  bhojpurOnAppInstallChange();
});

window.addEventListener('appinstalled', () => {
  deferredPrompt = null;
  bhojpurOnAppInstallChange();
});

function bhojpurIsAppInstallable() {
  return !bhojpurIsAppInstalled() && deferredPrompt != null;
}

function bhojpurIsAppInstalled() {
  const isStandalone = window.matchMedia('(display-mode: standalone)').matches;
  return isStandalone || navigator.standalone;
}

async function bhojpurShowInstallPrompt() {
  deferredPrompt.prompt();
  await deferredPrompt.userChoice;
  deferredPrompt = null;
}

// -----------------------------------------------------------------------------
// Keep body clean
// -----------------------------------------------------------------------------
function bhojpurKeepBodyClean() {
  const body = document.body;
  const bodyChildrenCount = body.children.length;

  const mutationObserver = new MutationObserver(function (mutationList) {
    mutationList.forEach((mutation) => {
      switch (mutation.type) {
        case 'childList':
          while (body.children.length > bodyChildrenCount) {
            body.removeChild(body.lastChild);
          }
          break;
      }
    });
  });

  mutationObserver.observe(document.body, {
    childList: true,
  });

  return () => mutationObserver.disconnect();
}

// -----------------------------------------------------------------------------
// Init Web Assembly
// -----------------------------------------------------------------------------
if (!/bot|googlebot|crawler|spider|robot|crawling/i.test(navigator.userAgent)) {
  if (!WebAssembly.instantiateStreaming) {
    WebAssembly.instantiateStreaming = async (resp, importObject) => {
      const source = await (await resp).arrayBuffer();
      return await WebAssembly.instantiate(source, importObject);
    };
  }

  const go = new Go();

  WebAssembly.instantiateStreaming(fetch("{{.Wasm}}"), go.importObject)
    .then(result => {
      const loaderIcon = document.getElementById("app-wasm-loader-icon");
      loaderIcon.className = "bhojpur-logo";

      go.run(result.instance);
    })
    .catch(err => {
      const loaderIcon = document.getElementById("app-wasm-loader-icon");
      loaderIcon.className = "bhojpur-logo";

      const loaderLabel = document.getElementById("app-wasm-loader-label");
      loaderLabel.innerText = err;

      console.error("loading Bhojpur Web application wasm failed: " + err);
    });
} else {
  document.getElementById('app-wasm-loader').style.display = "none";
}