import * as functions from "firebase-functions";

// eslint-disable-next-line @typescript-eslint/no-var-requires
const {Webhook, MessageBuilder} = require("discord-webhook-node");

// const hook = new Webhook("https://discord.com/api/webhooks/872966356620955718/mVYTmjIcP7_V_PjscznxQNcL4YJ9rf0Ul6DjmABzUc-XR869Kk0jSSnlNuNTuHRHB-rg");

const IMAGE_URL = "https://chaseapp.tv/icon.png";

export const discordWebhook = functions.https.onRequest(
    async (request, response) => {
      if (request.body.webhook) {
        const webhook = request.body.webhook;
        const embed = new MessageBuilder()
            .setTitle(webhook.name)
            .setAuthor("Chase Notifier", IMAGE_URL, "https://github.com/chase-app")
            .setURL(webhook.url)
            .setColor("#00b0f4")
            .setThumbnail(IMAGE_URL)
            .setDescription(webhook.desc)
            .setImage(webhook.imageurl)
            .setFooter("LIVE", "https://firebasestorage.googleapis.com/v0/b/chaseapp-8459b.appspot.com/o/images%2Fdot.jpg?alt=media&token=e9d1be0f-801b-464e-8598-f81459b2ad95")
            .setTimestamp();

        if (webhook.hook !== undefined) {
          const hook = new Webhook(webhook.hook);
          hook.setUsername("ChaseApp");
          hook.setAvatar(IMAGE_URL);
          await hook.send(embed);
        } else {
          console.log("No webhook provided");
          response.status(400).send("No webhook provided");
        }
      }
      response.status(200).end();
    }
);
