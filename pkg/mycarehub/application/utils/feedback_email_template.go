package utils

// FeedbackNotificationEmail if the supports feedback email template for feedback
const FeedbackNotificationEmail = `
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <title>New Feedback</title>
  <style>
    @import url('https://fonts.googleapis.com/css2?family=Red+Hat+Display:wght@300;400;500;600;700&display=swap');
  </style>
</head>
<body
  style="width: 100%; background: rgba(208,199,222,0.47);margin: 0;padding: 15px 0;font-family: 'Red Hat Display', sans-serif;color: #455A64;">

<!--  wrapper div start-->
<div
  style="width: 600px;background: white;padding: 30px;border-top: 5px solid #7453a5;box-sizing: border-box;margin: 0 auto;">

  <!--    table start-->
  <table style="width: 100%;">
    <!--      header start-->
    <tr>
      <td><img style="width: 100px;" src="https://storage.googleapis.com/mycarehub-test/media/original_images/MyAfyaHubDevLogo.png" alt=""></td>
    </tr>
    <tr style="width: 100%;">
      <td style="width: 100%;height: 200px;background: rgba(102,51,153,0.02);padding: 15px;text-align: center;"><img
        style="width: 80%;" src="https://bewell.co.ke/image/converse.png" alt=""></td>
    </tr>
    <!--      header end-->

    <!--      title start-->
    <tr>
      <td><h2 style="color: #53C451;font-size: 28px;margin: 25px 0;">New Feedback</h2></td>
    </tr>
    <!--      title end-->

    <!--      content area start-->
    <tr>
      <td><p style="font-size: 18px;margin: 0;">Hello,</p></td>
    </tr>
    <tr>
    <td>
      <div style="height: 10px;"></div>
    </td>
    <tr>
      <td><p style="font-size: 18px;margin: 0;">You have received feedback from {{.User}}</p></td>
    </tr>
    <tr>
      <td>
        <div style="height: 10px;"></div>
      </td>
    </tr>
    <tr>
      <td><p style="font-size: 18px;margin: 0;line-height: 1.7;"><strong>Below is the feedback:</strong></p></td>
    </tr>
    <!--      content area end-->
    <tr>
      <td>
        <div style="height: 10px;"></div>
      </td>
    </tr>
    <!--      footer area start-->
    <tr>
      <td>
        <div>
          <div style="width: 50%;margin-bottom: 30px;">
            <strong style="display: block;margin-bottom: 5px;">Feedback Message</strong>
            <span>{{.Message}}</span>
          </div>
        </div>
        <div>
            <div style="width: 50%;margin-bottom: 30px;">
                <strong style="display: block;margin-bottom: 5px;">Requires Follow Up</strong>
            <span>{{.RequiresFollowUp}}</span>
        </div>
        </div>
      </td>
    </tr>

  </table>
  <!--    table end-->

</div>
<!--  wrapper div end-->

</body>
</html>
`
