import AWS from "aws-sdk"

export const uploadToS3Bucket = async (file) => {
  const s3 = new AWS.S3({
    accessKeyId: process.env.NEXT_APP_ACCESS_KEY,
    secretAccessKey: process.env.NEXT_APP_SECRET_KEY,
  });
  const params = {
    Bucket: "hostego-aws-bucket",
    Key: file.name,
    Body: file,
    ContentType: file.type,
  };

  try {
    const data = await s3.upload(params).promise();
    const image = data.Location;
    return image;
  } catch (error) {
    console.log(error);
  }
};
