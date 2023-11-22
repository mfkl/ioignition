(function() {
  "use strict";

  // const currentLocation = window.location;
  const document = window.document;
  const currentScript = document.currentScript;
  const url = currentScript.getAttribute('data-domain')

  console.log(url)
})();
