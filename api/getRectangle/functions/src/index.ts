import * as functions from "firebase-functions";
import convex from "@turf/convex";
import {AllGeoJSON} from "@turf/helpers";
import {coordAll} from "@turf/meta";
import centroid from "@turf/centroid";
import transformRotate from "@turf/transform-rotate";
import bearing from "@turf/bearing";
import envelope from "@turf/envelope";
import area from "@turf/area";

// eslint-disable-next-line max-len
export const smallestSurroundingRectangleByArea = functions.https.onRequest((request, response) => {
  const geoJsonInput = request.body as unknown as AllGeoJSON;
  // eslint-disable-next-line @typescript-eslint/ban-ts-comment
  // @ts-ignore
  console.log("geoJsonInput", (geoJsonInput.geometry.coordinates));
  const convexHull = convex(geoJsonInput);
  // eslint-disable-next-line @typescript-eslint/ban-ts-comment
  // @ts-ignore
  const centroidCoords = centroid(convexHull);
  // eslint-disable-next-line @typescript-eslint/ban-ts-comment
  // @ts-ignore
  const allHullCoords = coordAll(convexHull);

  let minArea = Number.MAX_SAFE_INTEGER;
  let resultPolygon = null;

  for (let index = 0; index < allHullCoords.length - 1; index++) {
    const angle = bearing(allHullCoords[index], allHullCoords[index + 1]);

    // eslint-disable-next-line @typescript-eslint/ban-ts-comment
    // @ts-ignore
    const rotatedHull = transformRotate(convexHull, -1.0 * angle, {
      pivot: centroidCoords,
    });

    const envelopeOfHull = envelope(rotatedHull);
    const envelopeArea = area(envelopeOfHull);
    if (envelopeArea < minArea) {
      minArea = envelopeArea;
      resultPolygon = transformRotate(envelopeOfHull, angle, {
        pivot: centroidCoords,
      });
    }
  }
  console.log("Result", resultPolygon);
  response.send(resultPolygon);
});

