package aws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecs"
)

type taskAPI interface {
	RunTask(input *ecs.RunTaskInput) (*ecs.RunTaskOutput, error)
}

type Task struct {
	taskClient taskAPI
}

func NewTask(s *session.Session) *Task {
	return &Task{
		taskClient: ecs.New(s),
	}
}

func (t *Task) RunTask(input *ecs.RunTaskInput) (*ecs.RunTaskOutput, error) {
	_, err := t.taskClient.RunTask(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case ecs.ErrCodeServerException:
				return nil, fmt.Errorf("%s: %s", ecs.ErrCodeServerException, aerr.Error())
			case ecs.ErrCodeClientException:
				return nil, fmt.Errorf("%s: %s", ecs.ErrCodeClientException, aerr.Error())
			case ecs.ErrCodeInvalidParameterException:
				return nil, fmt.Errorf("%s: %s", ecs.ErrCodeInvalidParameterException, aerr.Error())
			case ecs.ErrCodeClusterNotFoundException:
				return nil, fmt.Errorf("%s: %s", ecs.ErrCodeClusterNotFoundException, aerr.Error())
			case ecs.ErrCodeUnsupportedFeatureException:
				return nil, fmt.Errorf("%s: %s", ecs.ErrCodeUnsupportedFeatureException, aerr.Error())
			case ecs.ErrCodePlatformUnknownException:
				return nil, fmt.Errorf("%s: %s", ecs.ErrCodePlatformUnknownException, aerr.Error())
			case ecs.ErrCodePlatformTaskDefinitionIncompatibilityException:
				return nil, fmt.Errorf("%s: %s", ecs.ErrCodePlatformTaskDefinitionIncompatibilityException, aerr.Error())
			case ecs.ErrCodeAccessDeniedException:
				return nil, fmt.Errorf("%s: %s", ecs.ErrCodeAccessDeniedException, aerr.Error())
			case ecs.ErrCodeBlockedException:
				return nil, fmt.Errorf("%s: %s", ecs.ErrCodeBlockedException, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			return nil, err
		}
		return nil, err
	}

	return nil, nil
}
