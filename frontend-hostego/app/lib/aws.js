"use client";
import AWS from "aws-sdk";
import imageCompression from 'browser-image-compression';
// Oprocess.envnly log in development
if (.NODE_ENV !== "production") {
  console.log("AWS Environment Variables:", {
    region: process.env.NEXT_PUBLIC_AWS_REGION,
    accessKey: process.env.NEXT_PUBLIC_AWS_ACCESS_KEY,
    secretKey: process.env.NEXT_PUBLIC_AWS_SECRET_KEY,
    bucket: process.env.NEXT_PUBLIC_AWS_BUCKET_NAME,
  });
}

// Configure AWS SDK
const awsConfig = {
  region: process.env.NEXT_PUBLIC_AWS_REGION,
  credentials: {
    accessKeyId: process.env.NEXT_PUBLIC_AWS_ACCESS_KEY,
    secretAccessKey: process.env.NEXT_PUBLIC_AWS_SECRET_KEY,
  },
};

// Initialize AWS with configuration
AWS.config.update(awsConfig);

// Create S3 service object
const s3 = new AWS.S3();

// Image compression options
const compressionOptions = {
  maxSizeMB: 0.2,          // Max file size in MB
  maxWidthOrHeight: 1024,  // Max width/height in pixels
  useWebWorker: true,      // Use web worker for better performance
  fileType: 'image/webp'   // Convert to WebP for better compression
};

export const uploadToS3Bucket = async (file) => {
  try {
    let compressedFile = file;

    // Only compress if it's an image
    if (file.type.startsWith('image/')) {
      // Show compression progress in development
      if (process.env.NODE_ENV !== "production") {
        console.log('Original file size:', file.size / 1024 / 1024, 'MB');
      }

      // Compress the image
      compressedFile = await imageCompression(file, compressionOptions);

      if (process.env.NODE_ENV !== "production") {
        console.log('Compressed file size:', compressedFile.size / 1024 / 1024, 'MB');
      }
    }

    // Create unique file name
    const uniqueFileName = `hostego_img_${Date.now()}_${Math.random().toString(36).substring(7)}.webp`;

    // Set up upload parameters
    const params = {
      Bucket: process.env.NEXT_PUBLIC_AWS_BUCKET_NAME,
      Key: uniqueFileName,
      Body: compressedFile,
      ContentType: 'image/webp',
      CacheControl: 'max-age=31536000', // Cache for 1 year
    };

    // Upload to S3
    const { Location } = await s3.upload(params).promise();
    return Location;
  } catch (error) {
    // More detailed error logging in development
    if (process.env.NODE_ENV !== "production") {
      console.error("S3 Upload Error:", error);
    }
    throw new Error("Failed to upload file to S3");
  }
};
