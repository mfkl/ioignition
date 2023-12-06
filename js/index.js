"use strict";

const localstoreId = 'ioEventId'
const sessionId = generateUUID()
let cleanup

function generateUUID() {
  return ([1e7] + -1e3 + -4e3 + -8e3 + -1e11).replace(/[018]/g, c =>
    (c ^ crypto.getRandomValues(new Uint8Array(1))[0] & 15 >> c / 4).toString(16)
  );
}

function generateId() {
  const id = 'io-' + generateUUID()
  try {
    localStorage.setItem(localstoreId, id)
  } catch (err) { }

  return id
}

function getClientId() {
  let uId = ''
  try {
    uId = localStorage.getItem(localstoreId)
  } catch (error) { }

  if (uId) return uId

  return generateId()
}

function getBrowserAndPlatform() {
  const userAgent = navigator.userAgent;
  const platformInfo = navigator.platform;
  let browser, platform;

  // Detecting browser name
  if (userAgent.includes("Firefox")) {
    browser = "Mozilla Firefox";
  } else if (userAgent.includes("SamsungBrowser")) {
    browser = "Samsung Internet";
  } else if (userAgent.includes("Opera") || userAgent.includes("OPR")) {
    browser = "Opera";
  } else if (userAgent.includes("Trident")) {
    browser = "Internet Explorer";
  } else if (userAgent.includes("Edge")) {
    browser = "Edge";
  } else if (userAgent.includes("Chrome")) {
    browser = "Chrome";
  } else if (userAgent.includes("Safari")) {
    browser = "Safari";
  } else {
    browser = "unknown";
  }

  // Detecting platform (operating system)
  if (platformInfo.startsWith("Win")) {
    platform = "Windows";
  } else if (platformInfo.startsWith("Mac")) {
    platform = "MacOS";
  } else if (platformInfo.startsWith("Linux")) {
    platform = "Linux";
  } else if (platformInfo.startsWith("iPhone") || platformInfo.startsWith("iPad")) {
    platform = "iOS";
  } else if (platformInfo.startsWith("Android")) {
    platform = "Android";
  } else {
    platform = "unknown";
  }

  return { browser, platform };
}

function sendEvent(eventName, data) {
  const isLocalhost =
    /^localhost$|^127(?:\.[0-9]+){0,2}\.[0-9]+$|^(?:0*:)*?:?0*1$/.test(
      location.hostname
    ) || location.protocol === 'file:';

  if (isLocalhost) {
    // return ---> uncomment
  }

  const clientId = getClientId()
  let keepalive = false

  let url = `${data.apiHost}/event/${clientId}`
  let payload = {
    sessionId,
    event: eventName,
    ...data
  };

  if (eventName === "sessionend") {
    url = url + "/end"
    payload = { sessionId, event: eventName }
    // https://developer.mozilla.org/en-US/docs/Web/API/Navigator/sendBeacon
    keepalive = true
  } else if (eventName === "pagechange") {
    url = url + "/url"
    payload = { sessionId, event: eventName, url: data.url }
  }

  fetch(`${data.apiHost}/event/${clientId}`, {
    method: 'POST',
    body: JSON.stringify(payload),
    headers: new Headers().append('Content-Type', 'text/json'),
    keepalive,
  })
}

function config() {
  const { browser, platform } = getBrowserAndPlatform()

  return {
    url: location.href,
    domain: document.currentScript.getAttribute('data-domain'),
    referrer: document.referrer || "",
    deviceWidth: window.innerWidth,
    browser,
    platform,
    apiHost: 'http://localhost:8080/api',
  }
}

function trackEvent(eventName) {
  sendEvent(eventName, config());
};

function trackPageview() {
  trackEvent('pageview');
};

function trackPageChange() {
  trackEvent('pagechange');
};

function endPageView() {
  if (document.visibilityState === 'hidden') {
    trackEvent('sessionend')
    // cleanup
    cleanup()
  }
}

function enableTracking() {
  // Attach pushState and popState listeners
  const originalPushState = history.pushState;
  if (originalPushState) {
    history.pushState = function(data, title, url) {
      originalPushState.apply(this, [data, title, url]);

      trackPageChange();
    };
  }

  addEventListener('visibilitychange', endPageView)

  // Trigger first page view
  trackPageview();

  return function cleanup() {
    if (originalPushState) {
      history.pushState = originalPushState;
    }
    removeEventListener('visibilitychange', endPageView)
  };
};

cleanup = enableTracking()
