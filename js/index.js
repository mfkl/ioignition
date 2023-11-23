"use strict";

const localstoreId = 'ioEventId'
let startTime

function generateEventId() {
  const id = 'io-' + new Date().getTime() + '-' + Math.random().toString(36).substring(2, 9);
  try {
    localStorage.setItem(localstoreId, id)
  } catch(err) {}

  return id
}

function getEventId() {
  let uId = ''
  try {
    uId = localStorage.getItem(localstoreId)
  } catch (error) {}
  
  if (uId) return uId

  return generateEventId()
}

function sendEvent(eventName, data) {
  const isLocalhost =
    /^localhost$|^127(?:\.[0-9]+){0,2}\.[0-9]+$|^(?:0*:)*?:?0*1$/.test(
      location.hostname
    ) || location.protocol === 'file:';

  if (isLocalhost) {
    // return ---> uncomment
  }

  const eventId = getEventId()
  const sessionTime = Date.now() - startTime
  
  const payload = {
    n: eventName,
    u: data.url,
    d: data.domain,
    r: data.referrer,
    w: data.deviceWidth,
    a: data.userAgent,
    t: sessionTime
  };

  const req = new XMLHttpRequest();
  req.open('POST', `${data.apiHost}/api/event/${eventId}`, true);
  req.setRequestHeader('Content-Type', 'text/plain');
  req.send(JSON.stringify(payload));

  req.onreadystatechange = () => {
    if (req.readyState === 4) return;
    // send some diagnostics??
    cleanup()
  };
}

function config() {
  return {
    url: location.href,
    domain: location.hostname,
    referrer: document.referrer || null,
    deviceWidth: window.innerWidth,
    userAgent: window.navigator.userAgent,
    apiHost: 'https://ioignition.com/api',
  }
}

function trackEvent (eventName, eventData) {
  sendEvent(eventName, { ...config(), ...eventData });
};

function trackPageview (eventData) {
  trackEvent('pageview', eventData);
};

function startTimeListener() {
  startTime = Date.now()
}

function enableTracking() {
  window.addEventListener('DOMContentLoaded', startTimeListener)
  window.addEventListener('beforeunload', trackPageview)
}

function cleanup() {
  window.removeEventListener('DOMContentLoaded', startTimeListener)
  window.removeEventListener('beforeunload', trackPageview)
}

enableTracking()
