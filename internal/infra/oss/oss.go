package oss

import (
	"fmt"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/google/uuid"
)

type Client struct {
	bucket *oss.Bucket
}

func New(endpoint, ak, sk, bucket string) (*Client, error) {
	client, err := oss.New(endpoint, ak, sk)
	if err != nil {
		return nil, err
	}
	b, err := client.Bucket(bucket)
	if err != nil {
		return nil, err
	}
	return &Client{bucket: b}, nil
}

func (c *Client) PresignPut(suffix string) (url string, object string, err error) {
	object = fmt.Sprintf("vehicle/%s.%s", uuid.NewString(), suffix)
	url, err = c.bucket.SignURL(object, oss.HTTPPut, 600, // 10min
		oss.ContentType("image/"+suffix))
	return url, object, err
}
