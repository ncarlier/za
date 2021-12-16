/**
 * This ZerØ Analytics script is provided under the MIT License (MIT)
 *
 * Copyright (c) 2020 Nicolas Carlier

 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:

 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.

 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

(function() { 
  'use strict';

  let queue = window.za.q || [];
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
    "error": trackError,
    "event": trackEvent,
  };

  function create(tid, domain) {
    config.tid = tid;
    config.domain = domain;
  }

  function send() {
    // abort when "Do Not Track" is set
    if ('doNotTrack' in navigator && navigator.doNotTrack === "1") {
      return;
    }
    var args = [].slice.call(arguments);
    var c = args.shift();
    trackers[c].apply(this, args);
  }

  function toURLSearchParams(obj) {
    return Object.keys(obj).reduce(function(value, k) {
      value.append(k, obj[k]);
      return value;
    }, new URLSearchParams());
  }

  function getData(source) {
    const ds = source + 'Storage';
    if (!(ds in window)) {
      return null;
    }
    const key = 'za-' + source;
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
    const key = 'za-' + source;
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
    const el = document.getElementById('za-script');
    return el ? el.src.replace(/za(\.min)?\.js/, 'collect') : '';
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

  function writeBeaconImg(q) {
    // create image tracker
    let img = document.createElement('img');
    img.setAttribute('alt', '');
    img.setAttribute('aria-hidden', 'true');
    img.src = getTrackerUrl() + '?' + toURLSearchParams(q).toString();
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

  let _actualOnErrorHandler = null;
  function trackError() {
    _actualOnErrorHandler = window.onerror;
    window.onerror = function(msg, url, line, column, error) {
      const q = {
        tid: config.tid,
        t: 'exception',
        exm: msg,
        exl: line,
        exc: column,
        exu: url,
        exe: error,
        z: Date.now(), // Cache buster
      };
      writeBeaconImg(q);
      if (_actualOnErrorHandler) {
        return _actualOnErrorHandler.apply(this, arguments);
      }
      return false;
    };
  }
  
  function trackEvent(payload) {
    const q = {
      tid: config.tid,
      t: 'event',
      d: window.btoa(JSON.stringify(payload)),
      z: Date.now(), // Cache buster
    };
    writeBeaconImg(q);
  }

  function trackPageView(options = {top: false}) {
    console.log('options', options); 
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
    // get user language
    const ul = navigator.language;
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
      t: 'pageview',
      dp: dp,
      dh: dh,
      dr: dr,
      ul: ul,
      nv: isNewVisitor() ? 1 : 0,
      ns: isNewSession() ? 1 : 0,
      z: Date.now(), // Cache buster
    };
    if (options.top && 'sendBeacon' in navigator) {
      document.addEventListener('visibilitychange', () => {
        if (document.visibilityState === 'hidden') {
          q.top = Date.now() - q.z;
          navigator.sendBeacon(getTrackerUrl(), toURLSearchParams(q));
        }
      });
    } else {
      writeBeaconImg(q);
    }
  }

  // define za global function
  window.za = function() {
    var args = [].slice.call(arguments);
    var c = args.shift();
    commands[c].apply(this, args);
  };

  // process command pipeline
  queue.forEach((i) => window.za.apply(this, i));
})();
