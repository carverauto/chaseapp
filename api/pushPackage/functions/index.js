import functions from 'firebase-functions';
import fs from "fs";
import path from "path";
// import express from "express";
import crypto from "crypto";
import forge from "node-forge";
import archiver from "archiver";

// const app = express();
const PORT = 5555;

const __dirname = path.resolve();

const iconFiles = [
    "icon_16x16.png",
    "icon_16x16@2x.png",
    "icon_32x32.png",
    "icon_32x32@2x.png",
    "icon_128x128.png",
    "icon_128x128@2x.png",
];

const websiteJson = {
    websiteName: "ChaseApp",
    websitePushID: "web.tv.chaseapp",
    allowedDomains: ["https://chaseapp.tv"],
    urlFormatString: "https://chaseapp.tv/chase/%@",
    authenticationToken: "d0nger19ern33d16ch4ract3rzatL333ST",
    webServiceURL: "https://us-central1-chaseapp-8459b.cloudfunctions.net/pushPackage",
};

const p12Asn1 = forge.asn1.fromDer(fs.readFileSync(
    __dirname + "/certs/apple_push.p12", "binary"));
const p12 = forge.pkcs12.pkcs12FromAsn1(
    p12Asn1);

const certBags = p12.getBags({bagType: forge.pki.oids.certBag});
const certBag = certBags[forge.pki.oids.certBag];
// eslint-disable-next-line @typescript-eslint/ban-ts-comment
// @ts-ignore
const cert = certBag[0].cert;

const keyBags = p12.getBags({bagType: forge.pki.oids.pkcs8ShroudedKeyBag});
const keyBag = keyBags[forge.pki.oids.pkcs8ShroudedKeyBag];

// eslint-disable-next-line @typescript-eslint/ban-ts-comment
// @ts-ignore
const key = keyBag[0].key;

const intermediate = forge.pki.certificateFromPem(
    fs.readFileSync(__dirname + "/certs/intermediate.pem", "utf8"));

/*
app.listen(PORT, () => {
    console.log(`Server listening on port ${PORT}`);
});
 */

//export function pushPackage(req, res) {
export const pushPackage = functions.https.onRequest(async (req, res) => {
    //app.post("/:version/pushPackages/:websitePushId", async (req, res) => {
    if (!cert) {
        console.log("cert is null");

        res.sendStatus(500);
        return;
    }

    if (!key) {
        console.log("key is null");

        res.sendStatus(500);
        return;
    }

    // const iconSourceDir = "...";
    const iconSourceDir = __dirname + "/icon.iconset";

    res.attachment("pushpackage.zip");

    const archive = archiver("zip");

    archive.on("error", function (err) {
        res.status(500).send({error: err.message});
    });

    archive.on("warning", function (err) {
        if (err.code === "ENOENT") {
            console.log(`Archive warning ${err}`);
        } else {
            throw err;
        }
    });

    archive.on("end", function () {
        console.log("Archive wrote %d bytes", archive.pointer());
    });

    archive.pipe(res);

    archive.directory(iconSourceDir, "icon.iconset");

    const manifest = {};

    const readPromises = [];

    iconFiles.forEach((i) =>
        readPromises.push(
            new Promise((resolve, reject) => {
                const hash = crypto.createHash("sha512");
                const readStream = fs.createReadStream(
                    path.join(iconSourceDir, i),
                    {encoding: "utf8"}
                );

                readStream.on("data", (chunk) => {
                    hash.update(chunk);
                });

                readStream.on("end", () => {
                    const digest = hash.digest("hex");
                    manifest[`icon.iconset/${i}`] = {
                        hashType: "sha512",
                        hashValue: `${digest}`,
                    };
                    resolve();
                });

                readStream.on("error", (err) => {
                    console.log(`Error on readStream for ${i}; ${err}`);
                    // eslint-disable-next-line prefer-promise-reject-errors
                    reject();
                });
            })
        )
    );

    try {
        await Promise.all(readPromises);
    } catch (error) {
        console.log(`Error writing files; ${error}`);

        res.sendStatus(500);
        return;
    }

    const webJSON = {
        ...websiteJson,
        ...{authenticationToken: "..."},
    };
    const webHash = crypto.createHash("sha512");

    const webJSONString = JSON.stringify(webJSON);

    webHash.update(webJSONString);

    manifest["website.json"] = {
        hashType: "sha512",
        hashValue: `${webHash.digest("hex")}`,
    };

    const manifestJSONString = JSON.stringify(manifest);

    archive.append(webJSONString, {name: "website.json"});
    archive.append(manifestJSONString, {name: "manifest.json"});

    const p7 = forge.pkcs7.createSignedData();
    p7.content = forge.util.createBuffer(manifestJSONString, "utf8");
    p7.addCertificate(cert);
    p7.addCertificate(intermediate);
    p7.addSigner({
        // eslint-disable-next-line @typescript-eslint/ban-ts-comment
        // @ts-ignore
        key,
        certificate: cert,
        digestAlgorithm: forge.pki.oids.sha256,
        authenticatedAttributes: [{
            type: forge.pki.oids.contentType,
            value: forge.pki.oids.data,
        }, {
            type: forge.pki.oids.messageDigest,
        }, {
            type: forge.pki.oids.signingTime,
            value: new Date().toString(),
        }],
    });
    p7.sign({detached: true});

    const pem = forge.pkcs7.messageToPem(p7);
    archive.append(Buffer.from(pem, "binary"), {name: "signature"});

    // Have also tried this:
    // archive.append(forge.asn1.toDer(p7.toAsn1()).getBytes(),
    // { name: "signature" });

    try {
        await archive.finalize();
    } catch (error) {
        console.log(`Error on archive.finalize(); ${error}`);

        res.sendStatus(500);
        return;
    }
});