"use strict";

const localstoreId = 'ioEventId'
let startTime

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
    url: data.url,
    domain: data.domain,
    referrer: data.referrer,
    width: data.deviceWidth,
    agent: data.userAgent,
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
  return {
    url: location.href,
    domain: document.currentScript.getAttribute('data-domain'),
    referrer: document.referrer || "",
    deviceWidth: window.innerWidth,
    userAgent: window.navigator.userAgent,
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
