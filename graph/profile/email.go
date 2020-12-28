package profile

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"strings"
	"text/template"

	"gitlab.slade360emr.com/go/base"
)

func getTester(ctx context.Context, email string) (*TesterWhitelist, error) {
	filter := &base.FilterInput{
		FilterBy: []*base.FilterParam{
			{
				FieldName:           "email",
				FieldType:           base.FieldTypeString,
				ComparisonOperation: base.OperationEqual,
				FieldValue:          email,
			},
		},
	}
	docs, _, err := base.QueryNodes(ctx, nil, filter, nil, &TesterWhitelist{})
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve tester with email %s: %s", email, err)
	}
	if len(docs) == 0 {
		return nil, nil
	}
	if len(docs) > 1 {
		return nil, fmt.Errorf("there is more than one tester record with the specified email")
	}

	tDoc := docs[0]
	tester := &TesterWhitelist{}
	err = tDoc.DataTo(tester)
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshal TesterWhitelist from firestore doc: %w", err)
	}
	return tester, nil
}

func isTester(ctx context.Context, emails []string) bool {
	isTester := false
	for _, email := range emails {
		if strings.Contains(email, "apple.com") {
			isTester = true // special case for Apple.com app reviewers
		}
		tester, err := getTester(ctx, email)
		if err != nil {
			log.Printf("unable to retrieve tester with email %s: %s", email, err)
			continue
		}
		if tester != nil {
			isTester = true
		}
	}
	return isTester
}

func generateProcessKYCApprovalEmailTemplate() string {
	t := template.Must(template.New("approvalKYCEmail").Parse(processKYCApprovalEmail))
	buf := new(bytes.Buffer)
	err := t.Execute(buf, "")
	if err != nil {
		log.Fatalf("Error while generating KYC approval email template: %s", err)
	}
	return buf.String()
}

func generateProcessKYCRejectionEmailTemplate() string {
	t := template.Must(template.New("rejectionKYCEmail").Parse(processKYCRejectionEmail))
	buf := new(bytes.Buffer)
	err := t.Execute(buf, "")
	if err != nil {
		log.Fatalf("Error while generating KYC rejection email template: %s", err)
	}
	return buf.String()
}
