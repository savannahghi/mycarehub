"""
Phone browser OS detection.

This is a simple HTTP python function that detects a devices OS
through their browser and redirect them to Playstore if they are on android
or AppStore if the are on iOS.
"""

import flask
from user_agents import parse

app = flask.Flask(__name__)

IOS = "iOS"
ANDROID = "Android"
A_LANDING_PAGE = "https://a.bewell.co.ke"
PLAY_STORE_LINK = "https://appdistribution.firebase.dev/i/eb0b5d95e67a3b3f"
APPLE_STORE_LINK = "https://testflight.apple.com/join/p2GAbpaz"


events = {
    "android": "redirected_to_android_playstore",
    "IOS": "redirected_to_iOS_appstore",
}


def htmlTemplate(event, link):
    """Format and return a html template."""
    return f""" 
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Be.Well By Slade360</title>
    <!-- Google Tag Manager -->
    <script>(function(w,d,s,l,i){{w[l]=w[l]||[];w[l].push({{'gtm.start':
    new Date().getTime(),event:'gtm.js'}});var f=d.getElementsByTagName(s)[0],
    j=d.createElement(s),dl=l!='dataLayer'?'&l='+l:'';j.async=true;j.src=
    'https://www.googletagmanager.com/gtm.js?id='+i+dl;f.parentNode.insertBefore(j,f);
    }})(window,document,'script','dataLayer','GTM-T5V349R');</script>
    <!-- End Google Tag Manager -->
</head>
<body>
    <!-- Google Tag Manager (noscript) -->
    <noscript><iframe src="https://www.googletagmanager.com/ns.html?id=GTM-T5V349R"
    height="0" width="0" style="display:none;visibility:hidden"></iframe></noscript>
    <!-- End Google Tag Manager (noscript) -->
    
    <!-- Start of HubSpot Embed Code -->
    <script type="text/javascript" id="hs-script-loader" async defer src="//js.hs-scripts.com/20198195.js"></script>
    <!-- End of HubSpot Embed Code -->

    <!-- The core Firebase JS SDK  -->
    <script src="https://www.gstatic.com/firebasejs/8.7.0/firebase-app.js"></script>
    <script src="https://www.gstatic.com/firebasejs/8.7.0/firebase-analytics.js"></script>
     <!-- AdRoll tracking pixel -->
    <script type="text/javascript">
      adroll_adv_id = "G34MBP2POFA2VKFCJHXOZ4";
      adroll_pix_id = "IMUEF2EGPBCOLDIDDDSZMU";
      adroll_vaersion = "2.0";
      (function (w, d, e, o, a) {{
        w.__adroll_loaded = true;
        w.adroll = w.adroll || [];
        w.adroll.f = ['setProperties', 'identify', 'track'];
        var roundtripUrl = "https://s.adroll.com/j/" + adroll_adv_id + "/roundtrip.js";
        for (a = 0; a < w.adroll.f.length; a++) {{
          w.adroll[w.adroll.f[a]] = w.adroll[w.adroll.f[a]] || (function (n) {{
            return function () {{ w.adroll.push([n, arguments]) }}
          }})(w.adroll.f[a])
        }} e = d.createElement('script');
        o = d.getElementsByTagName('script')[0];
        e.async = 1;
        e.src = roundtripUrl;
        o.parentNode.insertBefore(e, o);
      }})(window, document); adroll.track("pageView");
    </script>

    <!-- Facebook Pixel Code -->
    <script>
        !function (f, b, e, v, n, t, s) {{
            if (f.fbq) return; n = f.fbq = function () {{
                n.callMethod ?
                n.callMethod.apply(n, arguments) : n.queue.push(arguments)
            }};
            if (!f._fbq) f._fbq = n; n.push = n; n.loaded = !0; n.version = '2.0';
            n.queue = []; t = b.createElement(e); t.async = !0;
            t.src = v; s = b.getElementsByTagName(e)[0];
            s.parentNode.insertBefore(t, s)
        }}(window, document, 'script',
        'https://connect.facebook.net/en_US/fbevents.js');
        fbq('init', '400335678066977');
        fbq('track', 'PageView');
    </script>
    <noscript>
        <img height="1" width="1" style="display:none"
        src="https://www.facebook.com/tr?id=400335678066977&ev=PageView&noscript=1" />
    </noscript>
    <!-- End Facebook Pixel Code -->

</body>
</html>
"""  # noqa


def detect_browser(request):
    """
    Detect a browser's user-agent.

    Given the family of OS we get, we redirect to either
    our Play store or App store.
    """
    user_agent = parse(request.headers.get("User-Agent"))
    os_family = user_agent.os.family

    if os_family == IOS:
        return flask.render_template_string(
            htmlTemplate(events["IOS"], APPLE_STORE_LINK)
        )

    if os_family == ANDROID:
        return flask.render_template_string(
            htmlTemplate(events["android"], PLAY_STORE_LINK)
        )

    else:
        return flask.redirect(A_LANDING_PAGE)


@app.route("/")
def index():
    """Flask app entrypoint."""
    return detect_browser(flask.request)
