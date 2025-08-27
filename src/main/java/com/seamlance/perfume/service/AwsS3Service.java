package com.seamlance.perfume.service;

import com.amazonaws.auth.AWSStaticCredentialsProvider;
import com.amazonaws.auth.BasicAWSCredentials;
import com.amazonaws.regions.Regions;
import com.amazonaws.services.s3.AmazonS3;
import com.amazonaws.services.s3.AmazonS3ClientBuilder;
import com.amazonaws.services.s3.model.ObjectMetadata;
import com.amazonaws.services.s3.model.PutObjectRequest;
import com.seamlance.perfume.constants.Constant;
import com.seamlance.perfume.constants.ErrorConstant;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Service;
import org.springframework.web.multipart.MultipartFile;

import java.io.IOException;
import java.io.InputStream;

@Service
public class AwsS3Service {
    @Value("${aws.s3.bucket.name}")
    private String bucketName;

    @Value("${aws.s3.access.key}")
    private String awsAccessKey;

    @Value("${aws.s3.secret.key}")
    private String awsSecretKey;



    public String saveImageToS3(MultipartFile file) {
        String s3LocationImage = null;

        try {
            String s3Filename = file.getOriginalFilename();

            BasicAWSCredentials awsCredentials = new BasicAWSCredentials(awsAccessKey, awsSecretKey);
            AmazonS3 s3Client = AmazonS3ClientBuilder.standard()
                    .withCredentials(new AWSStaticCredentialsProvider(awsCredentials))
                    .withRegion(Regions.DEFAULT_REGION)
                    .build();

            InputStream inputStream = file.getInputStream();
            ObjectMetadata metadata = new ObjectMetadata();
            metadata.setContentType(Constant.IMAGE_FILE_CONTENT_TYPE);

            PutObjectRequest putObjectRequest = new PutObjectRequest(bucketName, s3Filename, inputStream, metadata);
            s3Client.putObject(putObjectRequest);

            s3LocationImage = "https://" + bucketName + ".s3.amazonawas.com/"+s3Filename;
        } catch (IOException e) {
            throw new RuntimeException(ErrorConstant.UNABLE_TO_UPLOAD_FILE + ": " + e.getMessage());
        }

        return s3LocationImage;
    }
}
