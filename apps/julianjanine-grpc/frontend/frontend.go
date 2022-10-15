package frontend

import (
	"fmt"

	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/cloudfront"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/route53"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/s3"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/ssm"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

func CreateFrontend(ctx *pulumi.Context) error {
	conf := config.New(ctx, "")
	frontendDomainName := conf.Require("frontend_domain_name")
	ssmCertificateKey := conf.Require("ssm_certificate_key")
	route53ZoneName := conf.Require("route53_zone_name")

	policy := fmt.Sprintf(`
  {
      "Version": "2012-10-17",
      "Statement": [
          {
              "Sid": "PublicRead",
              "Effect": "Allow",
              "Principal": "*",
              "Action": "s3:GetObject",
              "Resource": "arn:aws:s3:::%s/*"
          }
      ]
  }`, frontendDomainName)

	// creates s3 bucket
	bucket, err := s3.NewBucket(ctx, "frontend_bucket", &s3.BucketArgs{
		Acl:    pulumi.String("public-read"),
		Bucket: pulumi.String(frontendDomainName),
		Policy: pulumi.String(policy),
		Website: s3.BucketWebsiteArgs{
			IndexDocument: pulumi.String("index.html"),
			ErrorDocument: pulumi.String("index.html"),
		},
	})
	if err != nil {
		return err
	}

	certificateArn, err := ssm.LookupParameter(ctx, &ssm.LookupParameterArgs{
		Name: ssmCertificateKey,
	}, nil)
	if err != nil {
		return err
	}

	// create cloudfront distribution
	s3OriginId := getS3OriginId(bucket.ID())

	distribution, err := cloudfront.NewDistribution(ctx, "frontend_distribution", &cloudfront.DistributionArgs{
		Origins: cloudfront.DistributionOriginArray{
			&cloudfront.DistributionOriginArgs{
				DomainName: bucket.BucketRegionalDomainName,
				OriginId:   s3OriginId,
			},
		},

		Enabled:           pulumi.Bool(true),
		IsIpv6Enabled:     pulumi.Bool(true),
		DefaultRootObject: pulumi.String("index.html"),

		Aliases: pulumi.StringArray{
			pulumi.String(frontendDomainName),
		},

		Restrictions: cloudfront.DistributionRestrictionsArgs{
			GeoRestriction: cloudfront.DistributionRestrictionsGeoRestrictionArgs{
				RestrictionType: pulumi.String("none"),
			},
		},

		// Default Cache behaviuour
		DefaultCacheBehavior: cloudfront.DistributionDefaultCacheBehaviorArgs{
			AllowedMethods: pulumi.StringArray{
				pulumi.String("GET"),
				pulumi.String("HEAD"),
			},
			CachedMethods: pulumi.StringArray{
				pulumi.String("GET"),
				pulumi.String("HEAD"),
			},
			TargetOriginId:       s3OriginId,
			MinTtl:               pulumi.IntPtr(0),
			DefaultTtl:           pulumi.IntPtr(3600),
			MaxTtl:               pulumi.IntPtr(86400),
			ViewerProtocolPolicy: pulumi.String("allow-all"),
			ForwardedValues: cloudfront.DistributionDefaultCacheBehaviorForwardedValuesArgs{
				QueryString: pulumi.Bool(false),

				Cookies: cloudfront.DistributionDefaultCacheBehaviorForwardedValuesCookiesArgs{
					Forward: pulumi.String("none"),
				},
			},
		},

		// Cache behaviour with precedence 0
		OrderedCacheBehaviors: cloudfront.DistributionOrderedCacheBehaviorArray{
			&cloudfront.DistributionOrderedCacheBehaviorArgs{
				PathPattern: pulumi.String("*"),
				AllowedMethods: pulumi.StringArray{
					pulumi.String("GET"),
					pulumi.String("HEAD"),
				},
				CachedMethods: pulumi.StringArray{
					pulumi.String("GET"),
					pulumi.String("HEAD"),
				},
				TargetOriginId:       s3OriginId,
				MinTtl:               pulumi.IntPtr(0),
				DefaultTtl:           pulumi.IntPtr(86400),
				MaxTtl:               pulumi.IntPtr(31536000),
				ViewerProtocolPolicy: pulumi.String("redirect-to-https"),
				ForwardedValues: cloudfront.DistributionOrderedCacheBehaviorForwardedValuesArgs{
					QueryString: pulumi.Bool(false),
					Cookies: cloudfront.DistributionOrderedCacheBehaviorForwardedValuesCookiesArgs{
						Forward: pulumi.String("none"),
					},
				},
			},
		},

		// Custom error response to handle 4xx
		CustomErrorResponses: cloudfront.DistributionCustomErrorResponseArray{
			&cloudfront.DistributionCustomErrorResponseArgs{
				ErrorCode:          pulumi.Int(403),
				ErrorCachingMinTtl: pulumi.Int(60),
				ResponseCode:       pulumi.Int(200),
				ResponsePagePath:   pulumi.String("/index.html"),
			},
			&cloudfront.DistributionCustomErrorResponseArgs{
				ErrorCode:          pulumi.Int(404),
				ErrorCachingMinTtl: pulumi.Int(60),
				ResponseCode:       pulumi.Int(200),
				ResponsePagePath:   pulumi.String("/index.html"),
			},
		},

		// Certificate
		ViewerCertificate: cloudfront.DistributionViewerCertificateArgs{
			AcmCertificateArn:      pulumi.String(certificateArn.Value),
			SslSupportMethod:       pulumi.String("sni-only"),
			MinimumProtocolVersion: pulumi.String("TLSv1.2_2021"),
		},
	})
	if err != nil {
		return err
	}

	// Create route53 record
	zone, err := route53.LookupZone(ctx, &route53.LookupZoneArgs{
		Name:        pulumi.StringRef(route53ZoneName),
		PrivateZone: pulumi.BoolRef(false),
	}, nil)
	if err != nil {
		return err
	}

	_, err = route53.NewRecord(ctx, "frontend_route_53", &route53.RecordArgs{
		ZoneId: pulumi.String(zone.ZoneId),
		Name:   pulumi.String(""),
		Type:   pulumi.String("A"),
		Aliases: route53.RecordAliasArray{
			&route53.RecordAliasArgs{
				Name:                 distribution.DomainName,
				ZoneId:               distribution.HostedZoneId,
				EvaluateTargetHealth: pulumi.Bool(false),
			},
		},
	})
	if err != nil {
		return err
	}

	return nil
}

func getS3OriginId(input pulumi.StringInput) pulumi.StringOutput {
	return input.ToStringOutput().ApplyT(func(s string) string {
		return fmt.Sprintf("S3-%s", s)
	}).(pulumi.StringOutput)
}
