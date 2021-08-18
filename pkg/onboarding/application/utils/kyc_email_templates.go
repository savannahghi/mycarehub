package utils

// ProcessKYCApprovalEmail ...
const ProcessKYCApprovalEmail = `
<!DOCTYPE html>
<html lang="en" xmlns="http://www.w3.org/1999/xhtml" xmlns:v="urn:schemas-microsoft-com:vml"
    xmlns:o="urn:schemas-microsoft-com:office:office">

<head>
    <title>Be.Well Professional by Slade 360° - Connected healthcare platform. </title>
    <meta property="description" content="KYC details approved">
    <!--VIEWPORT-->
    <meta name="viewport" content="width=device-width; initial-scale=1.0; maximum-scale=1.0; user-scalable=no;">
    <meta name="viewport" content="width=600, initial-scale = 2.3, user-scalable=no">
    <meta name="viewport" content="width=device-width">
    <!--CHARSET-->
    <meta http-equiv="Content-Type" content="text/html; charset=utf-8" />

    <!-- IE=Edge and IE=X -->
    <meta http-equiv="X-UA-Compatible" content="IE=7" />
    <meta http-equiv="X-UA-Compatible" content="IE=8" />
    <meta http-equiv="X-UA-Compatible" content="IE=9" />
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <!--[if !mso]>-->
    <meta http-equiv="X-UA-Compatible" content="IE=edge" />
    <!--<![endif]-->

    <link rel="preconnect" href="https://fonts.googleapis.com">
    <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
    <link href="https://fonts.googleapis.com/css2?family=Red+Hat+Display:wght@400;500;700;900&display=swap"
        rel="stylesheet">

    <!-- CSS Reset : BEGIN -->
    <style>
        html,
        body {
            margin: 0 auto !important;
            padding: 0 !important;
            height: 100% !important;
            width: 100% !important;
            background: #f1f1f1;
        }

        /* What it does: Stops email clients resizing small text. */
        * {
            -ms-text-size-adjust: 100%;
            -webkit-text-size-adjust: 100%;
        }

        /* What it does: Centers email on Android 4.4 */
        div[style*="margin: 16px 0"] {
            margin: 0 !important;
        }

        /* What it does: Stops Outlook from adding extra spacing to tables. */
        table,

        /* What it does: Fixes webkit padding issue. */
        table {
            border-spacing: 0 !important;
            border-collapse: collapse !important;
            table-layout: fixed !important;
            margin: 0 auto !important;
        }

        /* What it does: Uses a better rendering method when resizing images in IE. */
        img {
            -ms-interpolation-mode: bicubic;
        }

        /* What it does: Prevents Windows 10 Mail from underlining links despite inline CSS. Styles for underlined links should be inline. */
        a {
            text-decoration: none;
        }

        /* What it does: A work-around for email clients meddling in triggered links. */
        *[x-apple-data-detectors],
        /* iOS */
        .unstyle-auto-detected-links *,
        .aBn {
            border-bottom: 0 !important;
            cursor: default !important;
            color: inherit !important;
            text-decoration: none !important;
            font-size: inherit !important;
            font-family: inherit !important;
            font-weight: inherit !important;
            line-height: inherit !important;
        }

        /* What it does: Prevents Gmail from displaying a download button on large, non-linked images. */
        .a6S {
            display: none !important;
            opacity: 0.01 !important;
        }

        /* What it does: Prevents Gmail from changing the text color in conversation threads. */
        .im {
            color: inherit !important;
        }

        /* If the above doesn't work, add a .g-img class to any image in question. */
        img.g-img+div {
            display: none !important;
        }

        /* What it does: Removes right gutter in Gmail iOS app: https://github.com/TedGoas/Cerberus/issues/89  */
        /* Create one of these media queries for each additional viewport size you'd like to fix */

        /* iPhone 4, 4S, 5, 5S, 5C, and 5SE */
        @media only screen and (min-device-width: 320px) and (max-device-width: 374px) {
            u~div .email-container {
                min-width: 320px !important;
            }
        }

        /* iPhone 6, 6S, 7, 8, and X */
        @media only screen and (min-device-width: 375px) and (max-device-width: 413px) {
            u~div .email-container {
                min-width: 375px !important;
            }
        }

        /* iPhone 6+, 7+, and 8+ */
        @media only screen and (min-device-width: 414px) {
            u~div .email-container {
                min-width: 414px !important;
            }
        }
    </style>

    <!-- CSS Reset : END -->

    <!-- Progressive Enhancements : BEGIN -->
    <style>
        .bg_white {
            background: #ffffff;
        }

        .bg_light {
            background: #fafafa;
        }

        .bg_purple {
            background: #7B54C4;
        }

        .email-section {
            padding: 2.5em;
        }

        h1,
        h2,
        h3,
        h4,
        h5,
        h6 {
            font-family: 'Red Hat Display', sans-serif;
            color: #000000;
            margin-top: 0;
            font-weight: 400;
        }

        body {
            font-family: 'Red Hat Display', sans-serif;
            font-weight: 400;
            font-size: 15px;
            line-height: 1.8;
            color: rgba(0, 0, 0, .4);
        }

        a {
            color: #2f89fc;
        }

        /*LOGO*/

        .logo h1 {
            margin: 0;
        }

        .logo h1 a {
            color: #000000;
            font-size: 20px;
            font-weight: 700;
            text-transform: uppercase;
            font-family: 'Red Hat Display', sans-serif;
        }

        p {
            color: #000000;
            font-size: 16px;
        }

        /*FOOTER*/

        .footer {
            color: rgba(255, 255, 255, .5);

        }

        .footer .heading {
            color: #ffffff;
            font-size: 14px;
        }

        .footer ul {
            margin: 0;
            padding: 0;
        }

        .footer ul li {
            list-style: none;
            margin-bottom: 16px;
            font-size: 12px;
            font-weight: 700;
        }

        h3 .footer-text {
            color: #f2f2f2;
        }

        .footer ul li a {
            color: rgba(255, 255, 255, 1);
        }


        @media screen and (max-width: 500px) {}
    </style>
</head>

<body width="100%" style="margin: 0; padding: 0 !important; background-color: #f1f1f1;">
    <center style="width: 100%; background-color: #f1f1f1;">
        <div
            style="display: none; font-size: 1px;max-height: 0px; max-width: 0px; opacity: 0; overflow: hidden; font-family: 'Red Hat Display', sans-serif;">
            &zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;
        </div>
        <div style="max-width: 600px; margin: 0 auto;" class="email-container">
            <!-- BEGIN BODY -->
            <table align="center" role="presentation" cellspacing="0" cellpadding="0" border="0" width="100%"
                style="margin: auto;">
                <tr>
                    <td valign="top" class="bg_white" style="padding: 1em 2.5em;">
                        <table role="presentation" border="0" cellpadding="0" cellspacing="0" width="100%">
                            <tr>
                                <td bgcolor="#ffffff" align="center" valign="top" style="
                                padding: 40px 20px 10px 20px;
                                ">
                                    <img src="https://lh3.googleusercontent.com/pw/ACtC-3fN_p8U8EZgmtQymnwrhr_-5Go6Kw5e5U7lkjyk1jjMIEwSs6rDNELplpgVk2IciMfw5AbnphxJYwdocnsE6Y88xyKGlNXm1E1x3Sm9uxeMHhwjf8YgNwo622G8cb-d7ntlbNl7-uPCEylu5O_KzZY=s638-no"
                                        width="125" height="120"
                                        style="display: block; border: 0px; margin-bottom: 0" />
                                </td>
                            </tr>
                        </table>
                    </td>
                </tr><!-- end tr -->
                <tr>
                    <td valign="middle" class="hero hero-2 bg_white" style="padding: 4em 0;">
                        <table>
                            <tr>
                                <td>
                                    <div class="text" style="padding: 0 3em;">
                                        <p style="margin-bottom: 14px">Hello,</p>
                                        <p style="margin: 0">
                                            Your KYC details have been reviewed and approved. We look forward
                                            to working with you.
                                        </p>
                                        <p style="margin-top: 24px">
                                            Regards,<br />
                                            Bev from Be.Well
                                        </p>
                                    </div>
                                </td>
                            </tr>
                        </table>
                    </td>
                </tr><!-- end tr -->
                <!-- 1 Column Text + Button : END -->
            </table>
            <table align="center" role="presentation" cellspacing="0" cellpadding="0" border="0" width="100%"
                style="margin: auto;">
                <tr>
                    <td valign="middle" class="bg_purple footer email-section">
                        <table>
                            <tr>
                                <td valign="top" width="33.333%" style="padding-top: 20px;">
                                    <table role="presentation" cellspacing="0" cellpadding="0" border="0" width="100%">
                                        <tr>
                                            <td style="text-align: left; padding-left: 10px;">
                                                <h3 class="heading">Social links</h3>
                                                <ul>
                                                    <li><a href="https://twitter.com/BeWellApp_">Twitter</a></li>
                                                    <li><a href="https://www.facebook.com/BeWellbySlade360">Facebook</a>
                                                    </li>
                                                    <li><a
                                                            href="https://www.instagram.com/BeWellBySlade360/">Instagram</a>
                                                    </li>
                                                    <li><a
                                                            href="https://www.linkedin.com/showcase/be-well-by-slade-360">Linkedin</a>
                                                </ul>
                                            </td>
                                        </tr>
                                    </table>
                                </td>
                                <td valign="top" width="33.333%" style="padding-top: 20px;">
                                    <table role="presentation" cellspacing="0" cellpadding="0" border="0" width="100%">
                                        <tr>
                                            <td style="text-align: left; padding-left: 10px;">
                                                <h3 class="heading">Company</h3>
                                                <ul>
                                                    <li><a href="https://bewell.co.ke/privacy.html">Privacy</a></li>
                                                </ul>
                                            </td>
                                        </tr>
                                    </table>
                                </td>
                                <td valign="top" width="33.333%" style="padding-top: 20px;">
                                    <table role="presentation" cellspacing="0" cellpadding="0" border="0" width="100%">
                                        <tr>
                                            <td style="text-align: left; padding-left: 5px; padding-right: 5px;">
                                                <h3 class="heading">Contact Info</h3>
                                                <ul>
                                                    <li><span class="text">For more information or queries, contact us
                                                            via</span></li>
                                                    <li><span class="text"><a href="tel:+254 790 360 360">+254 790 360
                                                                360</span></a></li>
                                                    <li><span class="text"><a
                                                                href="mailto:feedback@bewell.co.ke">feedback@bewell.co.ke</span></a>
                                                    </li>
                                                </ul>
                                            </td>
                                        </tr>
                                    </table>
                                </td>
                            </tr>
                        </table>
                    </td>
                </tr><!-- end: tr -->
                <tr>
                    <td valign="middle" class="bg_purple footer email-section">
                        <table>
                            <tr>
                                <td valign="top" width="33.333%">
                                    <table role="presentation" cellspacing="0" cellpadding="0" border="0" width="100%">
                                        <tr>
                                            <td style=" text-align: left; padding-right: 10px; font-size: 12px;">
                                                <span class="footer-text"></span>&copy; Savannah Informatics
                                                Limited.</span>
                                            </td>
                                        </tr>
                                    </table>
                                </td>
                            </tr>
                        </table>
                    </td>
                </tr>
            </table>

        </div>
    </center>
    <script src="https://cdn.jsdelivr.net/npm/publicalbum@latest/embed-ui.min.js" async></script>

    <script src="https://www.gstatic.com/firebasejs/8.7.1/firebase-app.js"></script>

    <script src="https://www.gstatic.com/firebasejs/8.7.1/firebase-analytics.js"></script>

    <script>
        var firebaseConfig = {
            apiKey: "AIzaSyAv2aRsSSHkOR6xGwwaw6-UTkvED3RNlBQ",
            authDomain: "bewell-app.firebaseapp.com",
            databaseURL: "https://bewell-app.firebaseio.com",
            projectId: "bewell-app",
            storageBucket: "bewell-app.appspot.com",
            messagingSenderId: "841947754847",
            appId: "1:841947754847:web:6304157d32c82fd96686ea",
            measurementId: "G-6XTZEB5070"
        };
        firebase.initializeApp(firebaseConfig);
        const analytics = firebase.analytics();

        analytics.logEvent('opened_processed_kyc_approved_email');
    </script>
</body>

</html>
`

// ProcessKYCRejectionEmail ...
const ProcessKYCRejectionEmail = `
<!DOCTYPE html>
<html lang="en" xmlns="http://www.w3.org/1999/xhtml" xmlns:v="urn:schemas-microsoft-com:vml"
    xmlns:o="urn:schemas-microsoft-com:office:office">

<head>
    <title>Be.Well Professional by Slade 360° - Connected healthcare platform. </title>
    <meta property="description" content="KYC detail rejected with reasons">
    <!--VIEWPORT-->
    <meta name="viewport" content="width=device-width; initial-scale=1.0; maximum-scale=1.0; user-scalable=no;">
    <meta name="viewport" content="width=600, initial-scale = 2.3, user-scalable=no">
    <meta name="viewport" content="width=device-width">
    <!--CHARSET-->
    <meta http-equiv="Content-Type" content="text/html; charset=utf-8" />

    <!-- IE=Edge and IE=X -->
    <meta http-equiv="X-UA-Compatible" content="IE=7" />
    <meta http-equiv="X-UA-Compatible" content="IE=8" />
    <meta http-equiv="X-UA-Compatible" content="IE=9" />
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <!--[if !mso]>-->
    <meta http-equiv="X-UA-Compatible" content="IE=edge" />
    <!--<![endif]-->

    <link rel="preconnect" href="https://fonts.googleapis.com">
    <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
    <link href="https://fonts.googleapis.com/css2?family=Red+Hat+Display:wght@400;500;700;900&display=swap"
        rel="stylesheet">

    <!-- CSS Reset : BEGIN -->
    <style>
        html,
        body {
            margin: 0 auto !important;
            padding: 0 !important;
            height: 100% !important;
            width: 100% !important;
            background: #f1f1f1;
        }

        /* What it does: Stops email clients resizing small text. */
        * {
            -ms-text-size-adjust: 100%;
            -webkit-text-size-adjust: 100%;
        }

        /* What it does: Centers email on Android 4.4 */
        div[style*="margin: 16px 0"] {
            margin: 0 !important;
        }

        /* What it does: Stops Outlook from adding extra spacing to tables. */
        table,

        /* What it does: Fixes webkit padding issue. */
        table {
            border-spacing: 0 !important;
            border-collapse: collapse !important;
            table-layout: fixed !important;
            margin: 0 auto !important;
        }

        /* What it does: Uses a better rendering method when resizing images in IE. */
        img {
            -ms-interpolation-mode: bicubic;
        }

        /* What it does: Prevents Windows 10 Mail from underlining links despite inline CSS. Styles for underlined links should be inline. */
        a {
            text-decoration: none;
        }

        /* What it does: A work-around for email clients meddling in triggered links. */
        *[x-apple-data-detectors],
        /* iOS */
        .unstyle-auto-detected-links *,
        .aBn {
            border-bottom: 0 !important;
            cursor: default !important;
            color: inherit !important;
            text-decoration: none !important;
            font-size: inherit !important;
            font-family: inherit !important;
            font-weight: inherit !important;
            line-height: inherit !important;
        }

        /* What it does: Prevents Gmail from displaying a download button on large, non-linked images. */
        .a6S {
            display: none !important;
            opacity: 0.01 !important;
        }

        /* What it does: Prevents Gmail from changing the text color in conversation threads. */
        .im {
            color: inherit !important;
        }

        /* If the above doesn't work, add a .g-img class to any image in question. */
        img.g-img+div {
            display: none !important;
        }

        /* What it does: Removes right gutter in Gmail iOS app: https://github.com/TedGoas/Cerberus/issues/89  */
        /* Create one of these media queries for each additional viewport size you'd like to fix */

        /* iPhone 4, 4S, 5, 5S, 5C, and 5SE */
        @media only screen and (min-device-width: 320px) and (max-device-width: 374px) {
            u~div .email-container {
                min-width: 320px !important;
            }
        }

        /* iPhone 6, 6S, 7, 8, and X */
        @media only screen and (min-device-width: 375px) and (max-device-width: 413px) {
            u~div .email-container {
                min-width: 375px !important;
            }
        }

        /* iPhone 6+, 7+, and 8+ */
        @media only screen and (min-device-width: 414px) {
            u~div .email-container {
                min-width: 414px !important;
            }
        }
    </style>

    <!-- CSS Reset : END -->

    <!-- Progressive Enhancements : BEGIN -->
    <style>
        .bg_white {
            background: #ffffff;
        }

        .bg_light {
            background: #fafafa;
        }

        .bg_purple {
            background: #7B54C4;
        }

        .email-section {
            padding: 2.5em;
        }

        h1,
        h2,
        h3,
        h4,
        h5,
        h6 {
            font-family: 'Red Hat Display', sans-serif;
            color: #000000;
            margin-top: 0;
            font-weight: 400;
        }

        body {
            font-family: 'Red Hat Display', sans-serif;
            font-weight: 400;
            font-size: 15px;
            line-height: 1.8;
            color: rgba(0, 0, 0, .4);
        }

        a {
            color: #2f89fc;
        }

        /*LOGO*/

        .logo h1 {
            margin: 0;
        }

        .logo h1 a {
            color: #000000;
            font-size: 20px;
            font-weight: 700;
            text-transform: uppercase;
            font-family: 'Red Hat Display', sans-serif;
        }

        p {
            color: #000000;
            font-size: 16px;
        }

        /*FOOTER*/

        .footer {
            color: rgba(255, 255, 255, .5);

        }

        .footer .heading {
            color: #ffffff;
            font-size: 14px;
        }

        .footer ul {
            margin: 0;
            padding: 0;
        }

        .footer ul li {
            list-style: none;
            margin-bottom: 16px;
            font-size: 12px;
            font-weight: 700;
        }

        h3 .footer-text {
            color: #f2f2f2;
        }

        .footer ul li a {
            color: rgba(255, 255, 255, 1);
        }


        @media screen and (max-width: 500px) {}
    </style>
</head>

<body width="100%" style="margin: 0; padding: 0 !important; background-color: #f1f1f1;">
    <center style="width: 100%; background-color: #f1f1f1;">
        <div
            style="display: none; font-size: 1px;max-height: 0px; max-width: 0px; opacity: 0; overflow: hidden; font-family: 'Red Hat Display', sans-serif;">
            &zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;
        </div>
        <div style="max-width: 600px; margin: 0 auto;" class="email-container">
            <!-- BEGIN BODY -->
            <table align="center" role="presentation" cellspacing="0" cellpadding="0" border="0" width="100%"
                style="margin: auto;">
                <tr>
                    <td valign="top" class="bg_white" style="padding: 1em 2.5em;">
                        <table role="presentation" border="0" cellpadding="0" cellspacing="0" width="100%">
                            <tr>
                                <td bgcolor="#ffffff" align="center" valign="top" style="
                                padding: 40px 20px 10px 20px;
                                ">
                                    <img src="https://lh3.googleusercontent.com/pw/ACtC-3fN_p8U8EZgmtQymnwrhr_-5Go6Kw5e5U7lkjyk1jjMIEwSs6rDNELplpgVk2IciMfw5AbnphxJYwdocnsE6Y88xyKGlNXm1E1x3Sm9uxeMHhwjf8YgNwo622G8cb-d7ntlbNl7-uPCEylu5O_KzZY=s638-no"
                                        width="125" height="120"
                                        style="display: block; border: 0px; margin-bottom: 0" />
                                </td>
                            </tr>
                        </table>
                    </td>
                </tr><!-- end tr -->
                <tr>
                    <td valign="middle" class="hero hero-2 bg_white" style="padding: 4em 0;">
                        <table>
                            <tr>
                                <td>
                                    <div class="text" style="padding: 0 3em;">
                                        <p style="margin-bottom: 14px">Hello,</p>
                                        <p style="margin: 0">
                                            Your KYC details have been reviewed and unfortunately not approved
                                            because of the following:
                                        </p>
                                        <p></p>
                                        <p>{{.Reason}}</p>
                                        <p></p>
                                        <p style="margin: 0">
                                            If you believe this was a mistake, please contact us via <br>
                                            <a href="tel:0790360360">+254 790 360 360</a> in order to resolve the issue.
                                        </p>

                                        <p style="margin-top: 24px">
                                            Regards,<br />
                                            Bowi from Be.Well
                                        </p>
                                    </div>
                                </td>
                            </tr>
                        </table>
                    </td>
                </tr><!-- end tr -->
                <!-- 1 Column Text + Button : END -->
            </table>
            <table align="center" role="presentation" cellspacing="0" cellpadding="0" border="0" width="100%"
                style="margin: auto;">
                <tr>
                    <td valign="middle" class="bg_purple footer email-section">
                        <table>
                            <tr>
                                <td valign="top" width="33.333%" style="padding-top: 20px;">
                                    <table role="presentation" cellspacing="0" cellpadding="0" border="0" width="100%">
                                        <tr>
                                            <td style="text-align: left; padding-left: 10px;">
                                                <h3 class="heading">Social links</h3>
                                                <ul>
                                                    <li><a href="https://twitter.com/BeWellApp_">Twitter</a></li>
                                                    <li><a href="https://www.facebook.com/BeWellbySlade360">Facebook</a>
                                                    </li>
                                                    <li><a
                                                            href="https://www.instagram.com/BeWellBySlade360/">Instagram</a>
                                                    </li>
                                                    <li><a
                                                            href="https://www.linkedin.com/showcase/be-well-by-slade-360">Linkedin</a>
                                                </ul>
                                            </td>
                                        </tr>
                                    </table>
                                </td>
                                <td valign="top" width="33.333%" style="padding-top: 20px;">
                                    <table role="presentation" cellspacing="0" cellpadding="0" border="0" width="100%">
                                        <tr>
                                            <td style="text-align: left; padding-left: 10px;">
                                                <h3 class="heading">Company</h3>
                                                <ul>
                                                    <li><a href="https://bewell.co.ke/privacy.html">Privacy</a></li>
                                                </ul>
                                            </td>
                                        </tr>
                                    </table>
                                </td>
                                <td valign="top" width="33.333%" style="padding-top: 20px;">
                                    <table role="presentation" cellspacing="0" cellpadding="0" border="0" width="100%">
                                        <tr>
                                            <td style="text-align: left; padding-left: 5px; padding-right: 5px;">
                                                <h3 class="heading">Contact Info</h3>
                                                <ul>
                                                    <li><span class="text">For more information or queries, contact us
                                                            via</span></li>
                                                    <li><span class="text"><a href="tel:+254 790 360 360">+254 790 360
                                                                360</span></a></li>
                                                    <li><span class="text"><a
                                                                href="mailto:feedback@bewell.co.ke">feedback@bewell.co.ke</span></a>
                                                    </li>
                                                </ul>
                                            </td>
                                        </tr>
                                    </table>
                                </td>
                            </tr>
                        </table>
                    </td>
                </tr><!-- end: tr -->
                <tr>
                    <td valign="middle" class="bg_purple footer email-section">
                        <table>
                            <tr>
                                <td valign="top" width="33.333%">
                                    <table role="presentation" cellspacing="0" cellpadding="0" border="0" width="100%">
                                        <tr>
                                            <td style=" text-align: left; padding-right: 10px; font-size: 12px;">
                                                <span class="footer-text"></span>&copy; Savannah Informatics
                                                Limited.</span>
                                            </td>
                                        </tr>
                                    </table>
                                </td>
                            </tr>
                        </table>
                    </td>
                </tr>
            </table>

        </div>
    </center>
    <script src="https://cdn.jsdelivr.net/npm/publicalbum@latest/embed-ui.min.js" async></script>
    <!-- The core Firebase JS SDK is always required and must be listed first -->
    <script src="https://www.gstatic.com/firebasejs/8.7.1/firebase-app.js"></script>

    <script src="https://www.gstatic.com/firebasejs/8.7.1/firebase-analytics.js"></script>

    <script>
        var firebaseConfig = {
            apiKey: "AIzaSyAv2aRsSSHkOR6xGwwaw6-UTkvED3RNlBQ",
            authDomain: "bewell-app.firebaseapp.com",
            databaseURL: "https://bewell-app.firebaseio.com",
            projectId: "bewell-app",
            storageBucket: "bewell-app.appspot.com",
            messagingSenderId: "841947754847",
            appId: "1:841947754847:web:6304157d32c82fd96686ea",
            measurementId: "G-6XTZEB5070"
        };
        firebase.initializeApp(firebaseConfig);
        const analytics = firebase.analytics();

        analytics.logEvent('opened_processed_kyc_rejection_email');
    </script>
</body>

</html>
`

// AcknowledgementKYCEmail ...
const AcknowledgementKYCEmail = `
<!DOCTYPE html>
<html lang="en" xmlns="http://www.w3.org/1999/xhtml" xmlns:v="urn:schemas-microsoft-com:vml"
    xmlns:o="urn:schemas-microsoft-com:office:office">

<head>
    <title>Be.Well Professional by Slade 360° - Connected healthcare platform. </title>
    <meta property="description" content="Acknowledged receipt of KYC details.">
    <!--VIEWPORT-->
    <meta name="viewport" content="width=device-width; initial-scale=1.0; maximum-scale=1.0; user-scalable=no;">
    <meta name="viewport" content="width=600, initial-scale = 2.3, user-scalable=no">
    <meta name="viewport" content="width=device-width">
    <!--CHARSET-->
    <meta http-equiv="Content-Type" content="text/html; charset=utf-8" />

    <!-- IE=Edge and IE=X -->
    <meta http-equiv="X-UA-Compatible" content="IE=7" />
    <meta http-equiv="X-UA-Compatible" content="IE=8" />
    <meta http-equiv="X-UA-Compatible" content="IE=9" />
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <!--[if !mso]>-->
    <meta http-equiv="X-UA-Compatible" content="IE=edge" />
    <!--<![endif]-->

    <link rel="preconnect" href="https://fonts.googleapis.com">
    <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
    <link href="https://fonts.googleapis.com/css2?family=Red+Hat+Display:wght@400;500;700;900&display=swap"
        rel="stylesheet">

    <!-- CSS Reset : BEGIN -->
    <style>
        html,
        body {
            margin: 0 auto !important;
            padding: 0 !important;
            height: 100% !important;
            width: 100% !important;
            background: #f1f1f1;
        }

        /* What it does: Stops email clients resizing small text. */
        * {
            -ms-text-size-adjust: 100%;
            -webkit-text-size-adjust: 100%;
        }

        /* What it does: Centers email on Android 4.4 */
        div[style*="margin: 16px 0"] {
            margin: 0 !important;
        }

        /* What it does: Stops Outlook from adding extra spacing to tables. */
        table,

        /* What it does: Fixes webkit padding issue. */
        table {
            border-spacing: 0 !important;
            border-collapse: collapse !important;
            table-layout: fixed !important;
            margin: 0 auto !important;
        }

        /* What it does: Uses a better rendering method when resizing images in IE. */
        img {
            -ms-interpolation-mode: bicubic;
        }

        /* What it does: Prevents Windows 10 Mail from underlining links despite inline CSS. Styles for underlined links should be inline. */
        a {
            text-decoration: none;
        }

        /* What it does: A work-around for email clients meddling in triggered links. */
        *[x-apple-data-detectors],
        /* iOS */
        .unstyle-auto-detected-links *,
        .aBn {
            border-bottom: 0 !important;
            cursor: default !important;
            color: inherit !important;
            text-decoration: none !important;
            font-size: inherit !important;
            font-family: inherit !important;
            font-weight: inherit !important;
            line-height: inherit !important;
        }

        /* What it does: Prevents Gmail from displaying a download button on large, non-linked images. */
        .a6S {
            display: none !important;
            opacity: 0.01 !important;
        }

        /* What it does: Prevents Gmail from changing the text color in conversation threads. */
        .im {
            color: inherit !important;
        }

        /* If the above doesn't work, add a .g-img class to any image in question. */
        img.g-img+div {
            display: none !important;
        }

        /* What it does: Removes right gutter in Gmail iOS app: https://github.com/TedGoas/Cerberus/issues/89  */
        /* Create one of these media queries for each additional viewport size you'd like to fix */

        /* iPhone 4, 4S, 5, 5S, 5C, and 5SE */
        @media only screen and (min-device-width: 320px) and (max-device-width: 374px) {
            u~div .email-container {
                min-width: 320px !important;
            }
        }

        /* iPhone 6, 6S, 7, 8, and X */
        @media only screen and (min-device-width: 375px) and (max-device-width: 413px) {
            u~div .email-container {
                min-width: 375px !important;
            }
        }

        /* iPhone 6+, 7+, and 8+ */
        @media only screen and (min-device-width: 414px) {
            u~div .email-container {
                min-width: 414px !important;
            }
        }
    </style>

    <!-- CSS Reset : END -->

    <!-- Progressive Enhancements : BEGIN -->
    <style>
        .bg_white {
            background: #ffffff;
        }

        .bg_light {
            background: #fafafa;
        }

        .bg_purple {
            background: #7B54C4;
        }

        .email-section {
            padding: 2.5em;
        }

        h1,
        h2,
        h3,
        h4,
        h5,
        h6 {
            font-family: 'Red Hat Display', sans-serif;
            color: #000000;
            margin-top: 0;
            font-weight: 400;
        }

        body {
            font-family: 'Red Hat Display', sans-serif;
            font-weight: 400;
            font-size: 15px;
            line-height: 1.8;
            color: rgba(0, 0, 0, .4);
        }

        a {
            color: #2f89fc;
        }

        /*LOGO*/

        .logo h1 {
            margin: 0;
        }

        .logo h1 a {
            color: #000000;
            font-size: 20px;
            font-weight: 700;
            text-transform: uppercase;
            font-family: 'Red Hat Display', sans-serif;
        }

        p {
            color: #000000;
            font-size: 16px;
        }

        /*FOOTER*/

        .footer {
            color: rgba(255, 255, 255, .5);

        }

        .footer .heading {
            color: #ffffff;
            font-size: 14px;
        }

        .footer ul {
            margin: 0;
            padding: 0;
        }

        .footer ul li {
            list-style: none;
            margin-bottom: 16px;
            font-size: 12px;
            font-weight: 700;
        }

        h3 .footer-text {
            color: #f2f2f2;
        }

        .footer ul li a {
            color: rgba(255, 255, 255, 1);
        }


        @media screen and (max-width: 500px) {}
    </style>
</head>

<body width="100%" style="margin: 0; padding: 0 !important; background-color: #f1f1f1;">
    <center style="width: 100%; background-color: #f1f1f1;">
        <div
            style="display: none; font-size: 1px;max-height: 0px; max-width: 0px; opacity: 0; overflow: hidden; font-family: 'Red Hat Display', sans-serif;">
            &zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;
        </div>
        <div style="max-width: 600px; margin: 0 auto;" class="email-container">
            <!-- BEGIN BODY -->
            <table align="center" role="presentation" cellspacing="0" cellpadding="0" border="0" width="100%"
                style="margin: auto;">
                <tr>
                    <td valign="top" class="bg_white" style="padding: 1em 2.5em;">
                        <table role="presentation" border="0" cellpadding="0" cellspacing="0" width="100%">
                            <tr>
                                <td bgcolor="#ffffff" align="center" valign="top" style="
                                padding: 40px 20px 10px 20px;
                                ">
                                    <img src="https://lh3.googleusercontent.com/pw/ACtC-3fN_p8U8EZgmtQymnwrhr_-5Go6Kw5e5U7lkjyk1jjMIEwSs6rDNELplpgVk2IciMfw5AbnphxJYwdocnsE6Y88xyKGlNXm1E1x3Sm9uxeMHhwjf8YgNwo622G8cb-d7ntlbNl7-uPCEylu5O_KzZY=s638-no"
                                        width="125" height="120"
                                        style="display: block; border: 0px; margin-bottom: 0" />
                                </td>
                            </tr>
                        </table>
                    </td>
                </tr><!-- end tr -->
                <tr>
                    <td valign="middle" class="hero hero-2 bg_white" style="padding: 4em 0;">
                        <table>
                            <tr>
                                <td>
                                    <div class="text" style="padding: 0 3em;">
                                        <p style="margin-bottom: 14px">Dear {{.SupplierName}},</p>
                                        <p style="margin: 0">
                                            We have received your {{.AccountType}} {{.PartnerType}} KYC document.
                                        </p>
                                        <p style="margin: 0">
                                            We will review the request and we will be in touch.

                                        <p style="margin: 0">
                                            Thank you for using Be.Well.
                                        </p>

                                        <p style="margin-top: 24px">
                                            Regards,<br />
                                            Bev from Be.Well
                                        </p>
                                    </div>
                                </td>
                            </tr>
                        </table>
                    </td>
                </tr><!-- end tr -->
                <!-- 1 Column Text + Button : END -->
            </table>
            <table align="center" role="presentation" cellspacing="0" cellpadding="0" border="0" width="100%"
                style="margin: auto;">
                <tr>
                    <td valign="middle" class="bg_purple footer email-section">
                        <table>
                            <tr>
                                <td valign="top" width="33.333%" style="padding-top: 20px;">
                                    <table role="presentation" cellspacing="0" cellpadding="0" border="0" width="100%">
                                        <tr>
                                            <td style="text-align: left; padding-left: 10px;">
                                                <h3 class="heading">Social links</h3>
                                                <ul>
                                                    <li><a href="https://twitter.com/BeWellApp_">Twitter</a></li>
                                                    <li><a href="https://www.facebook.com/BeWellbySlade360">Facebook</a>
                                                    </li>
                                                    <li><a
                                                            href="https://www.instagram.com/BeWellBySlade360/">Instagram</a>
                                                    </li>
                                                    <li><a
                                                            href="https://www.linkedin.com/showcase/be-well-by-slade-360">Linkedin</a>
                                                </ul>
                                            </td>
                                        </tr>
                                    </table>
                                </td>
                                <td valign="top" width="33.333%" style="padding-top: 20px;">
                                    <table role="presentation" cellspacing="0" cellpadding="0" border="0" width="100%">
                                        <tr>
                                            <td style="text-align: left; padding-left: 10px;">
                                                <h3 class="heading">Company</h3>
                                                <ul>
                                                    <li><a href="https://bewell.co.ke/privacy.html">Privacy</a></li>
                                                </ul>
                                            </td>
                                        </tr>
                                    </table>
                                </td>
                                <td valign="top" width="33.333%" style="padding-top: 20px;">
                                    <table role="presentation" cellspacing="0" cellpadding="0" border="0" width="100%">
                                        <tr>
                                            <td style="text-align: left; padding-left: 5px; padding-right: 5px;">
                                                <h3 class="heading">Contact Info</h3>
                                                <ul>
                                                    <li><span class="text">For more information or queries, contact us
                                                            via</span></li>
                                                    <li><span class="text"><a href="tel:+254 790 360 360">+254 790 360
                                                                360</span></a></li>
                                                    <li><span class="text"><a
                                                                href="mailto:feedback@bewell.co.ke">feedback@bewell.co.ke</span></a>
                                                    </li>
                                                </ul>
                                            </td>
                                        </tr>
                                    </table>
                                </td>
                            </tr>
                        </table>
                    </td>
                </tr><!-- end: tr -->
                <tr>
                    <td valign="middle" class="bg_purple footer email-section">
                        <table>
                            <tr>
                                <td valign="top" width="33.333%">
                                    <table role="presentation" cellspacing="0" cellpadding="0" border="0" width="100%">
                                        <tr>
                                            <td style=" text-align: left; padding-right: 10px; font-size: 12px;">
                                                <span class="footer-text"></span>&copy; Savannah Informatics
                                                Limited.</span>
                                            </td>
                                        </tr>
                                    </table>
                                </td>
                            </tr>
                        </table>
                    </td>
                </tr>
            </table>

        </div>
    </center>
    <script src="https://cdn.jsdelivr.net/npm/publicalbum@latest/embed-ui.min.js" async></script>
    <!-- The core Firebase JS SDK is always required and must be listed first -->
    <script src="https://www.gstatic.com/firebasejs/8.7.1/firebase-app.js"></script>

    <script src="https://www.gstatic.com/firebasejs/8.7.1/firebase-analytics.js"></script>

    <script>
        var firebaseConfig = {
            apiKey: "AIzaSyAv2aRsSSHkOR6xGwwaw6-UTkvED3RNlBQ",
            authDomain: "bewell-app.firebaseapp.com",
            databaseURL: "https://bewell-app.firebaseio.com",
            projectId: "bewell-app",
            storageBucket: "bewell-app.appspot.com",
            messagingSenderId: "841947754847",
            appId: "1:841947754847:web:6304157d32c82fd96686ea",
            measurementId: "G-6XTZEB5070"
        };
        firebase.initializeApp(firebaseConfig);
        const analytics = firebase.analytics();

        analytics.logEvent('opened_acknowledgement_kyc_email');
    </script>
</body>

</html>
`

// AdminKYCSubmittedEmail ...
const AdminKYCSubmittedEmail = `
<!DOCTYPE html>
<html lang="en" xmlns="http://www.w3.org/1999/xhtml" xmlns:v="urn:schemas-microsoft-com:vml"
    xmlns:o="urn:schemas-microsoft-com:office:office">

<head>
    <title>Be.Well Professional by Slade 360° - Connected healthcare platform. </title>
    <meta property="description" content="Admin KYC details submitted.">
    <!--VIEWPORT-->
    <meta name="viewport" content="width=device-width; initial-scale=1.0; maximum-scale=1.0; user-scalable=no;">
    <meta name="viewport" content="width=600, initial-scale = 2.3, user-scalable=no">
    <meta name="viewport" content="width=device-width">
    <!--CHARSET-->
    <meta http-equiv="Content-Type" content="text/html; charset=utf-8" />

    <!-- IE=Edge and IE=X -->
    <meta http-equiv="X-UA-Compatible" content="IE=7" />
    <meta http-equiv="X-UA-Compatible" content="IE=8" />
    <meta http-equiv="X-UA-Compatible" content="IE=9" />
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <!--[if !mso]>-->
    <meta http-equiv="X-UA-Compatible" content="IE=edge" />
    <!--<![endif]-->

    <link rel="preconnect" href="https://fonts.googleapis.com">
    <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
    <link href="https://fonts.googleapis.com/css2?family=Red+Hat+Display:wght@400;500;700;900&display=swap"
        rel="stylesheet">

    <!-- CSS Reset : BEGIN -->
    <style>
        html,
        body {
            margin: 0 auto !important;
            padding: 0 !important;
            height: 100% !important;
            width: 100% !important;
            background: #f1f1f1;
        }

        /* What it does: Stops email clients resizing small text. */
        * {
            -ms-text-size-adjust: 100%;
            -webkit-text-size-adjust: 100%;
        }

        /* What it does: Centers email on Android 4.4 */
        div[style*="margin: 16px 0"] {
            margin: 0 !important;
        }

        /* What it does: Stops Outlook from adding extra spacing to tables. */
        table,

        /* What it does: Fixes webkit padding issue. */
        table {
            border-spacing: 0 !important;
            border-collapse: collapse !important;
            table-layout: fixed !important;
            margin: 0 auto !important;
        }

        /* What it does: Uses a better rendering method when resizing images in IE. */
        img {
            -ms-interpolation-mode: bicubic;
        }

        /* What it does: Prevents Windows 10 Mail from underlining links despite inline CSS. Styles for underlined links should be inline. */
        a {
            text-decoration: none;
        }

        /* What it does: A work-around for email clients meddling in triggered links. */
        *[x-apple-data-detectors],
        /* iOS */
        .unstyle-auto-detected-links *,
        .aBn {
            border-bottom: 0 !important;
            cursor: default !important;
            color: inherit !important;
            text-decoration: none !important;
            font-size: inherit !important;
            font-family: inherit !important;
            font-weight: inherit !important;
            line-height: inherit !important;
        }

        /* What it does: Prevents Gmail from displaying a download button on large, non-linked images. */
        .a6S {
            display: none !important;
            opacity: 0.01 !important;
        }

        /* What it does: Prevents Gmail from changing the text color in conversation threads. */
        .im {
            color: inherit !important;
        }

        /* If the above doesn't work, add a .g-img class to any image in question. */
        img.g-img+div {
            display: none !important;
        }

        /* What it does: Removes right gutter in Gmail iOS app: https://github.com/TedGoas/Cerberus/issues/89  */
        /* Create one of these media queries for each additional viewport size you'd like to fix */

        /* iPhone 4, 4S, 5, 5S, 5C, and 5SE */
        @media only screen and (min-device-width: 320px) and (max-device-width: 374px) {
            u~div .email-container {
                min-width: 320px !important;
            }
        }

        /* iPhone 6, 6S, 7, 8, and X */
        @media only screen and (min-device-width: 375px) and (max-device-width: 413px) {
            u~div .email-container {
                min-width: 375px !important;
            }
        }

        /* iPhone 6+, 7+, and 8+ */
        @media only screen and (min-device-width: 414px) {
            u~div .email-container {
                min-width: 414px !important;
            }
        }
    </style>

    <!-- CSS Reset : END -->

    <!-- Progressive Enhancements : BEGIN -->
    <style>
        .bg_white {
            background: #ffffff;
        }

        .bg_light {
            background: #fafafa;
        }

        .bg_purple {
            background: #7B54C4;
        }

        .email-section {
            padding: 2.5em;
        }

        h1,
        h2,
        h3,
        h4,
        h5,
        h6 {
            font-family: 'Red Hat Display', sans-serif;
            color: #000000;
            margin-top: 0;
            font-weight: 400;
        }

        body {
            font-family: 'Red Hat Display', sans-serif;
            font-weight: 400;
            font-size: 15px;
            line-height: 1.8;
            color: rgba(0, 0, 0, .4);
        }

        a {
            color: #2f89fc;
        }

        /*LOGO*/

        .logo h1 {
            margin: 0;
        }

        .logo h1 a {
            color: #000000;
            font-size: 20px;
            font-weight: 700;
            text-transform: uppercase;
            font-family: 'Red Hat Display', sans-serif;
        }

        p {
            color: #000000;
            font-size: 16px;
        }

        /*FOOTER*/

        .footer {
            color: rgba(255, 255, 255, .5);

        }

        .footer .heading {
            color: #ffffff;
            font-size: 14px;
        }

        .footer ul {
            margin: 0;
            padding: 0;
        }

        .footer ul li {
            list-style: none;
            margin-bottom: 16px;
            font-size: 12px;
            font-weight: 700;
        }

        h3 .footer-text {
            color: #f2f2f2;
        }

        .footer ul li a {
            color: rgba(255, 255, 255, 1);
        }


        @media screen and (max-width: 500px) {}
    </style>
</head>

<body width="100%" style="margin: 0; padding: 0 !important; background-color: #f1f1f1;">
    <center style="width: 100%; background-color: #f1f1f1;">
        <div
            style="display: none; font-size: 1px;max-height: 0px; max-width: 0px; opacity: 0; overflow: hidden; font-family: 'Red Hat Display', sans-serif;">
            &zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;
        </div>
        <div style="max-width: 600px; margin: 0 auto;" class="email-container">
            <!-- BEGIN BODY -->
            <table align="center" role="presentation" cellspacing="0" cellpadding="0" border="0" width="100%"
                style="margin: auto;">
                <tr>
                    <td valign="top" class="bg_white" style="padding: 1em 2.5em;">
                        <table role="presentation" border="0" cellpadding="0" cellspacing="0" width="100%">
                            <tr>
                                <td bgcolor="#ffffff" align="center" valign="top" style="
                                padding: 40px 20px 10px 20px;">
                                    <img src="https://lh3.googleusercontent.com/pw/ACtC-3fN_p8U8EZgmtQymnwrhr_-5Go6Kw5e5U7lkjyk1jjMIEwSs6rDNELplpgVk2IciMfw5AbnphxJYwdocnsE6Y88xyKGlNXm1E1x3Sm9uxeMHhwjf8YgNwo622G8cb-d7ntlbNl7-uPCEylu5O_KzZY=s638-no"
                                        width="125" height="120"
                                        style="display: block; border: 0px; margin-bottom: 0" />
                                </td>
                            </tr>
                        </table>
                    </td>
                </tr><!-- end tr -->
                <tr>
                    <td valign="middle" class="hero hero-2 bg_white" style="padding: 4em 0;">
                        <table>
                            <tr>
                                <td>
                                    <div class="text" style="padding: 0 3em; text-align: left;">
                                        <p style="margin-bottom: 14px;">Hello,</p>
                                        <p>{{.EmailBody}}</p>

                                        <p>Below are your supplier details:</p>
                                        <p>Partner Name: <span style="color: #000000;">{{.SupplierName}}</span>
                                        </p>
                                        <p style="margin: 0;">Partner Type: <span
                                                style="color: #000000;">{{.PartnerType}}</span>
                                        </p>
                                        <p style="margin: 0;">Account Type: <span
                                                style="color: #000000;">{{.AccountType}}</span>
                                        </p>
                                        <p style="margin: 0;">Email: <span
                                                style="color: #000000;">{{.EmailAddress}}</span></p>
                                        <p style="margin: 0;">Phone Number: <span
                                                style="color: #000000;">{{.PrimaryPhone}}</span>
                                        </p>

                                        <p style="margin-top: 24px">
                                            Regards,<br />
                                            Bev from Be.Well
                                        </p>
                                    </div>
                                </td>
                            </tr>
                        </table>
                    </td>
                </tr><!-- end tr -->
                <!-- 1 Column Text + Button : END -->
            </table>
            <table align="center" role="presentation" cellspacing="0" cellpadding="0" border="0" width="100%"
                style="margin: auto;">
                <tr>
                    <td valign="middle" class="bg_purple footer email-section">
                        <table>
                            <tr>
                                <td valign="top" width="33.333%" style="padding-top: 20px;">
                                    <table role="presentation" cellspacing="0" cellpadding="0" border="0" width="100%">
                                        <tr>
                                            <td style="text-align: left; padding-left: 10px;">
                                                <h3 class="heading">Social links</h3>
                                                <ul>
                                                    <li><a href="https://twitter.com/BeWellApp_">Twitter</a></li>
                                                    <li><a href="https://www.facebook.com/BeWellbySlade360">Facebook</a>
                                                    </li>
                                                    <li><a
                                                            href="https://www.instagram.com/BeWellBySlade360/">Instagram</a>
                                                    </li>
                                                    <li><a
                                                            href="https://www.linkedin.com/showcase/be-well-by-slade-360">Linkedin</a>
                                                </ul>
                                            </td>
                                        </tr>
                                    </table>
                                </td>
                                <td valign="top" width="33.333%" style="padding-top: 20px;">
                                    <table role="presentation" cellspacing="0" cellpadding="0" border="0" width="100%">
                                        <tr>
                                            <td style="text-align: left; padding-left: 10px;">
                                                <h3 class="heading">Company</h3>
                                                <ul>
                                                    <li><a href="https://bewell.co.ke/privacy.html">Privacy</a></li>
                                                </ul>
                                            </td>
                                        </tr>
                                    </table>
                                </td>
                                <td valign="top" width="33.333%" style="padding-top: 20px;">
                                    <table role="presentation" cellspacing="0" cellpadding="0" border="0" width="100%">
                                        <tr>
                                            <td style="text-align: left; padding-left: 5px; padding-right: 5px;">
                                                <h3 class="heading">Contact Info</h3>
                                                <ul>
                                                    <li><span class="text">For more information or queries, contact us
                                                            via</span></li>
                                                    <li><span class="text"><a href="tel:+254 790 360 360">+254 790 360
                                                                360</span></a></li>
                                                    <li><span class="text"><a
                                                                href="mailto:feedback@bewell.co.ke">feedback@bewell.co.ke</span></a>
                                                    </li>
                                                </ul>
                                            </td>
                                        </tr>
                                    </table>
                                </td>
                            </tr>
                        </table>
                    </td>
                </tr><!-- end: tr -->
                <tr>
                    <td valign="middle" class="bg_purple footer email-section">
                        <table>
                            <tr>
                                <td valign="top" width="33.333%">
                                    <table role="presentation" cellspacing="0" cellpadding="0" border="0" width="100%">
                                        <tr>
                                            <td style=" text-align: left; padding-right: 10px; font-size: 12px;">
                                                <span class="footer-text"></span>&copy; Savannah Informatics
                                                Limited.</span>
                                            </td>
                                        </tr>
                                    </table>
                                </td>
                            </tr>
                        </table>
                    </td>
                </tr>
            </table>

        </div>
    </center>
    <script src="https://cdn.jsdelivr.net/npm/publicalbum@latest/embed-ui.min.js" async></script>
    <!-- The core Firebase JS SDK is always required and must be listed first -->
    <script src="https://www.gstatic.com/firebasejs/8.7.1/firebase-app.js"></script>

    <script src="https://www.gstatic.com/firebasejs/8.7.1/firebase-analytics.js"></script>

    <script>
        var firebaseConfig = {
            apiKey: "AIzaSyAv2aRsSSHkOR6xGwwaw6-UTkvED3RNlBQ",
            authDomain: "bewell-app.firebaseapp.com",
            databaseURL: "https://bewell-app.firebaseio.com",
            projectId: "bewell-app",
            storageBucket: "bewell-app.appspot.com",
            messagingSenderId: "841947754847",
            appId: "1:841947754847:web:6304157d32c82fd96686ea",
            measurementId: "G-6XTZEB5070"
        };
        firebase.initializeApp(firebaseConfig);
        const analytics = firebase.analytics();

        analytics.logEvent('opened_admin_submitted_kyc_email');
    </script>
</body>

</html>
`

// SupplierSuspensionEmail ...
const SupplierSuspensionEmail = `
<!DOCTYPE html>
<html lang="en" xmlns="http://www.w3.org/1999/xhtml" xmlns:v="urn:schemas-microsoft-com:vml"
    xmlns:o="urn:schemas-microsoft-com:office:office">

<head>
    <title>Be.Well Professional by Slade 360° - Connected healthcare platform. </title>
    <meta property="description" content="Supplier has been suspended.">
    <!--VIEWPORT-->
    <meta name="viewport" content="width=device-width; initial-scale=1.0; maximum-scale=1.0; user-scalable=no;">
    <meta name="viewport" content="width=600, initial-scale = 2.3, user-scalable=no">
    <meta name="viewport" content="width=device-width">
    <!--CHARSET-->
    <meta http-equiv="Content-Type" content="text/html; charset=utf-8" />

    <!-- IE=Edge and IE=X -->
    <meta http-equiv="X-UA-Compatible" content="IE=7" />
    <meta http-equiv="X-UA-Compatible" content="IE=8" />
    <meta http-equiv="X-UA-Compatible" content="IE=9" />
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <!--[if !mso]>-->
    <meta http-equiv="X-UA-Compatible" content="IE=edge" />
    <!--<![endif]-->

    <link rel="preconnect" href="https://fonts.googleapis.com">
    <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
    <link href="https://fonts.googleapis.com/css2?family=Red+Hat+Display:wght@400;500;700;900&display=swap"
        rel="stylesheet">

    <!-- CSS Reset : BEGIN -->
    <style>
        html,
        body {
            margin: 0 auto !important;
            padding: 0 !important;
            height: 100% !important;
            width: 100% !important;
            background: #f1f1f1;
        }

        /* What it does: Stops email clients resizing small text. */
        * {
            -ms-text-size-adjust: 100%;
            -webkit-text-size-adjust: 100%;
        }

        /* What it does: Centers email on Android 4.4 */
        div[style*="margin: 16px 0"] {
            margin: 0 !important;
        }

        /* What it does: Stops Outlook from adding extra spacing to tables. */
        table,

        /* What it does: Fixes webkit padding issue. */
        table {
            border-spacing: 0 !important;
            border-collapse: collapse !important;
            table-layout: fixed !important;
            margin: 0 auto !important;
        }

        /* What it does: Uses a better rendering method when resizing images in IE. */
        img {
            -ms-interpolation-mode: bicubic;
        }

        /* What it does: Prevents Windows 10 Mail from underlining links despite inline CSS. Styles for underlined links should be inline. */
        a {
            text-decoration: none;
        }

        /* What it does: A work-around for email clients meddling in triggered links. */
        *[x-apple-data-detectors],
        /* iOS */
        .unstyle-auto-detected-links *,
        .aBn {
            border-bottom: 0 !important;
            cursor: default !important;
            color: inherit !important;
            text-decoration: none !important;
            font-size: inherit !important;
            font-family: inherit !important;
            font-weight: inherit !important;
            line-height: inherit !important;
        }

        /* What it does: Prevents Gmail from displaying a download button on large, non-linked images. */
        .a6S {
            display: none !important;
            opacity: 0.01 !important;
        }

        /* What it does: Prevents Gmail from changing the text color in conversation threads. */
        .im {
            color: inherit !important;
        }

        /* If the above doesn't work, add a .g-img class to any image in question. */
        img.g-img+div {
            display: none !important;
        }

        /* What it does: Removes right gutter in Gmail iOS app: https://github.com/TedGoas/Cerberus/issues/89  */
        /* Create one of these media queries for each additional viewport size you'd like to fix */

        /* iPhone 4, 4S, 5, 5S, 5C, and 5SE */
        @media only screen and (min-device-width: 320px) and (max-device-width: 374px) {
            u~div .email-container {
                min-width: 320px !important;
            }
        }

        /* iPhone 6, 6S, 7, 8, and X */
        @media only screen and (min-device-width: 375px) and (max-device-width: 413px) {
            u~div .email-container {
                min-width: 375px !important;
            }
        }

        /* iPhone 6+, 7+, and 8+ */
        @media only screen and (min-device-width: 414px) {
            u~div .email-container {
                min-width: 414px !important;
            }
        }
    </style>

    <!-- CSS Reset : END -->

    <!-- Progressive Enhancements : BEGIN -->
    <style>
        .bg_white {
            background: #ffffff;
        }

        .bg_light {
            background: #fafafa;
        }

        .bg_purple {
            background: #7B54C4;
        }

        .email-section {
            padding: 2.5em;
        }

        h1,
        h2,
        h3,
        h4,
        h5,
        h6 {
            font-family: 'Red Hat Display', sans-serif;
            color: #000000;
            margin-top: 0;
            font-weight: 400;
        }

        body {
            font-family: 'Red Hat Display', sans-serif;
            font-weight: 400;
            font-size: 15px;
            line-height: 1.8;
            color: rgba(0, 0, 0, .4);
        }

        a {
            color: #2f89fc;
        }

        /*LOGO*/

        .logo h1 {
            margin: 0;
        }

        .logo h1 a {
            color: #000000;
            font-size: 20px;
            font-weight: 700;
            text-transform: uppercase;
            font-family: 'Red Hat Display', sans-serif;
        }

        p {
            color: #000000;
            font-size: 16px;
        }

        /*FOOTER*/

        .footer {
            color: rgba(255, 255, 255, .5);

        }

        .footer .heading {
            color: #ffffff;
            font-size: 14px;
        }

        .footer ul {
            margin: 0;
            padding: 0;
        }

        .footer ul li {
            list-style: none;
            margin-bottom: 16px;
            font-size: 12px;
            font-weight: 700;
        }

        h3 .footer-text {
            color: #f2f2f2;
        }

        .footer ul li a {
            color: rgba(255, 255, 255, 1);
        }


        @media screen and (max-width: 500px) {}
    </style>
</head>

<body width="100%" style="margin: 0; padding: 0 !important; background-color: #f1f1f1;">
    <center style="width: 100%; background-color: #f1f1f1;">
        <div
            style="display: none; font-size: 1px;max-height: 0px; max-width: 0px; opacity: 0; overflow: hidden; font-family: 'Red Hat Display', sans-serif;">
            &zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;
        </div>
        <div style="max-width: 600px; margin: 0 auto;" class="email-container">
            <!-- BEGIN BODY -->
            <table align="center" role="presentation" cellspacing="0" cellpadding="0" border="0" width="100%"
                style="margin: auto;">
                <tr>
                    <td valign="top" class="bg_white" style="padding: 1em 2.5em;">
                        <table role="presentation" border="0" cellpadding="0" cellspacing="0" width="100%">
                            <tr>
                                <td bgcolor="#ffffff" align="center" valign="top" style="
                                padding: 40px 20px 10px 20px;
                                ">
                                    <img src="https://lh3.googleusercontent.com/pw/ACtC-3fN_p8U8EZgmtQymnwrhr_-5Go6Kw5e5U7lkjyk1jjMIEwSs6rDNELplpgVk2IciMfw5AbnphxJYwdocnsE6Y88xyKGlNXm1E1x3Sm9uxeMHhwjf8YgNwo622G8cb-d7ntlbNl7-uPCEylu5O_KzZY=s638-no"
                                        width="125" height="120"
                                        style="display: block; border: 0px; margin-bottom: 0" />
                                </td>
                            </tr>
                        </table>
                    </td>
                </tr><!-- end tr -->
                <tr>
                    <td valign="middle" class="hero hero-2 bg_white" style="padding: 4em 0;">
                        <table>
                            <tr>
                                <td>
                                    <div class="text" style="padding: 0 3em;">
                                        <p style="margin-bottom: 14px">Dear {{.SupplierName}},</p>
                                        <p style="margin: 0">
                                            {{.EmailBody}}
                                        </p>
                                        <p>You will not be able to transact on Be.Well while on suspension.</p>
                                        <p></p>
                                        <p>Incase of any queries, please contact us via <a href="tel:0790360360">+254
                                                790 360 360.</a></p>
                                        </p>

                                        <p style="margin-top: 24px">
                                            Regards,<br />
                                            Bev from Be.Well
                                        </p>
                                    </div>
                                </td>
                            </tr>
                        </table>
                    </td>
                </tr><!-- end tr -->
                <!-- 1 Column Text + Button : END -->
            </table>
            <table align="center" role="presentation" cellspacing="0" cellpadding="0" border="0" width="100%"
                style="margin: auto;">
                <tr>
                    <td valign="middle" class="bg_purple footer email-section">
                        <table>
                            <tr>
                                <td valign="top" width="33.333%" style="padding-top: 20px;">
                                    <table role="presentation" cellspacing="0" cellpadding="0" border="0" width="100%">
                                        <tr>
                                            <td style="text-align: left; padding-left: 10px;">
                                                <h3 class="heading">Social links</h3>
                                                <ul>
                                                    <li><a href="https://twitter.com/BeWellApp_">Twitter</a></li>
                                                    <li><a href="https://www.facebook.com/BeWellbySlade360">Facebook</a>
                                                    </li>
                                                    <li><a
                                                            href="https://www.instagram.com/BeWellBySlade360/">Instagram</a>
                                                    </li>
                                                    <li><a
                                                            href="https://www.linkedin.com/showcase/be-well-by-slade-360">Linkedin</a>
                                                </ul>
                                            </td>
                                        </tr>
                                    </table>
                                </td>
                                <td valign="top" width="33.333%" style="padding-top: 20px;">
                                    <table role="presentation" cellspacing="0" cellpadding="0" border="0" width="100%">
                                        <tr>
                                            <td style="text-align: left; padding-left: 10px;">
                                                <h3 class="heading">Company</h3>
                                                <ul>
                                                    <li><a href="https://bewell.co.ke/privacy.html">Privacy</a></li>
                                                </ul>
                                            </td>
                                        </tr>
                                    </table>
                                </td>
                                <td valign="top" width="33.333%" style="padding-top: 20px;">
                                    <table role="presentation" cellspacing="0" cellpadding="0" border="0" width="100%">
                                        <tr>
                                            <td style="text-align: left; padding-left: 5px; padding-right: 5px;">
                                                <h3 class="heading">Contact Info</h3>
                                                <ul>
                                                    <li><span class="text">For more information or queries, contact us
                                                            via</span></li>
                                                    <li><span class="text"><a href="tel:+254 790 360 360">+254 790 360
                                                                360</span></a></li>
                                                    <li><span class="text"><a
                                                                href="mailto:feedback@bewell.co.ke">feedback@bewell.co.ke</span></a>
                                                    </li>
                                                </ul>
                                            </td>
                                        </tr>
                                    </table>
                                </td>
                            </tr>
                        </table>
                    </td>
                </tr><!-- end: tr -->
                <tr>
                    <td valign="middle" class="bg_purple footer email-section">
                        <table>
                            <tr>
                                <td valign="top" width="33.333%">
                                    <table role="presentation" cellspacing="0" cellpadding="0" border="0" width="100%">
                                        <tr>
                                            <td style=" text-align: left; padding-right: 10px; font-size: 12px;">
                                                <span class="footer-text"></span>&copy; Savannah Informatics
                                                Limited.</span>
                                            </td>
                                        </tr>
                                    </table>
                                </td>
                            </tr>
                        </table>
                    </td>
                </tr>
            </table>

        </div>
    </center>
    <script src="https://cdn.jsdelivr.net/npm/publicalbum@latest/embed-ui.min.js" async></script>
    <!-- The core Firebase JS SDK is always required and must be listed first -->
    <script src="https://www.gstatic.com/firebasejs/8.7.1/firebase-app.js"></script>

    <script src="https://www.gstatic.com/firebasejs/8.7.1/firebase-analytics.js"></script>

    <script>
        var firebaseConfig = {
            apiKey: "AIzaSyAv2aRsSSHkOR6xGwwaw6-UTkvED3RNlBQ",
            authDomain: "bewell-app.firebaseapp.com",
            databaseURL: "https://bewell-app.firebaseio.com",
            projectId: "bewell-app",
            storageBucket: "bewell-app.appspot.com",
            messagingSenderId: "841947754847",
            appId: "1:841947754847:web:6304157d32c82fd96686ea",
            measurementId: "G-6XTZEB5070"
        };
        firebase.initializeApp(firebaseConfig);
        const analytics = firebase.analytics();


        analytics.logEvent('opened_supplier_suspension_email');
    </script>
</body>

</html>
`
