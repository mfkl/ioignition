"use strict";

const localstoreId = 'ioEventId'

function generateId() {
  const id = 'io-' + new Date().getTime() + '-' + Math.random().toString(36).substring(2, 9);
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
  // unique per instance
  const eventId = generateId()

  const payload = {
    eventid: eventId,
    event: eventName,
    ...data
  };

  const req = new XMLHttpRequest();
  req.open('POST', `${data.apiHost}/event/${clientId}`, true);
  req.setRequestHeader('Content-Type', 'text/json');
  req.send(JSON.stringify(payload));

  req.onreadystatechange = () => {
    if (req.readyState === 4) return;
    // send some diagnostics??
    if (eventName === 'sessionend') {
      cleanup()
    }
  };
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

function endPageView() {
  trackEvent('sessionend')
}

function enableTracking() {
  window.addEventListener('beforeunload', endPageView)
}

function cleanup() {
  window.removeEventListener('beforeunload', endPageView)
}

trackPageview()
enableTracking()
