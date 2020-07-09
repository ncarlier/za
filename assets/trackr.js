(function() { 
  'use strict';

  let queue = window.trackr.q || [];
  let config = {
    tid: '',
    domain: 'auto',
  };
  const commands = {
    "create": create,
    "send": send,
  };
  const trackers = {
    "pageview": trackPageView,
  };

  function create(tid, domain) {
    config.tid = tid;
    config.domain = domain;
  }

  function send() {
    var args = [].slice.call(arguments);
    var c = args.shift();
    trackers[c].apply(this, args);
  }

  function toQueryString(obj) {
    return '?' + Object.keys(obj).map(function(k) {
      return encodeURIComponent(k) + '=' + encodeURIComponent(obj[k]);
    }).join('&');
  }

  function getData(source) {
    const ds = source + 'Storage';
    if (!(ds in window)) {
      return null;
    }
    const key = 'trackr-' + source;
    const value = window[ds].getItem(key);
    if (!value) {
      return null;
    }
    try {
      return JSON.parse(value);
    } catch (err) {
      return null;
    }
  }

  function setData(source, data) {
    const ds = source + 'Storage';
    if (!(ds in window)) {
      return null;
    }
    const key = 'trackr-' + source;
    window[ds].setItem(key, JSON.stringify(data));
  }

  function isNewVisitor() {
    if (getData('local') === null) {
      setData('local', {firstSeen: Date.now()});
      return true;
    }
    return false;
  }

  function isNewSession() {
    if (getData('session') === null) {
      setData('session', {lastSeen: Date.now()});
      return true;
    }
    setData('session', {lastSeen: Date.now()});
    return false;
  }

  function getTrackerUrl() {
    const el = document.getElementById('trackr-script');
    return el ? el.src.replace('trackr.js', 'collect') : '';
  }

  function getCanonicalURL(loc) {
    let canonical = document.querySelector('link[rel="canonical"][href]');
    if (canonical) {
      let a = document.createElement('a');
      a.href = canonical.href;
      return a;
    }
    return loc;
  }

  function trackPageView() { 
    // abort when "Do Not Track" is set
    if ('doNotTrack' in navigator && navigator.doNotTrack === "1") {
      return;
    }
    // abort when page is in pre-rendered state
    if ('visibilityState' in document && document.visibilityState === 'prerender') {
      return;
    }
    // defer retry if page is not yet loaded
    if (document.body === null) {
      document.addEventListener("DOMContentLoaded", () => {
        trackPageView();
      });
      return;
    }
    // get document location
    let loc = window.location;
    // abort when the page is not provided by an HTTP server
    if (loc.host === '') {
      return;
    }
    // get canonical URL
    loc = getCanonicalURL(loc);
    // get document path
    let dp = loc.pathname + loc.search;
    if (!dp) {
      dp = '/';
    }
    // get document host name
    const dh = loc.protocol + "//" + loc.hostname;
    // get document referrer
    let dr = '';
    if (document.referrer.indexOf(dh) < 0) {
      dr = document.referrer;
    }
    // build the tracking query
    const q = {
      tid: config.tid,
      t: 'pageview', // TODO: 'pageview', 'screenview', 'event', 'transaction', 'item', 'social', 'exception', 'timing'
      dp: dp,
      dh: dh,
      dr: dr,
      nv: isNewVisitor() ? 1 : 0,
      ns: isNewSession() ? 1 : 0,
      z: Date.now(), // Cache buster
    };
    // create image tracker
    let img = document.createElement('img');
    img.setAttribute('alt', '');
    img.setAttribute('aria-hidden', 'true');
    img.src = getTrackerUrl() + toQueryString(q);
    img.addEventListener('load', function() {
      // remove image tracker from DOM
      document.body.removeChild(img);
    });
    
    // ensure to remove image tracker form DOM
    window.setTimeout(() => { 
      if (!img.parentNode) {
        return;
      }
      img.src = ''; 
      document.body.removeChild(img);
    }, 1000);
    // add image tracker to the DOM
    document.body.appendChild(img);  
  }

  // define Trackr global function
  window.trackr = function() {
    var args = [].slice.call(arguments);
    var c = args.shift();
    commands[c].apply(this, args);
  };

  // process command pipeline
  queue.forEach((i) => window.trackr.apply(this, i));
})();
