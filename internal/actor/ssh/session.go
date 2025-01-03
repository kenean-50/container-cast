package ssh

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/ssh"
)

type session struct {
	sess   *ssh.Session
	env    []string
	ctx    context.Context
	logger zerolog.Logger
}

func NewSession(client *ssh.Client) *session {
	if client == nil {
		log.Fatal().Msg("failed to create SSH session")
	}

	sess, err := client.NewSession()
	if err != nil {
		log.Fatal().Msg("failed to create SSH session")
	}
	// defer client.Close()

	return &session{
		sess: sess,
		ctx:  context.Background(),
		// env:  env,
		logger: log.
			With().
			Str("actor", "ssh").
			Logger(),
	}
}

// todo: set env for the session
// func (r *session) SetEnv() (err error) {
// 	for _, value := range r.env {
// 		env := strings.Split(value, "=")
// 		if err := r.sess.Setenv(env[0], strings.Join(env[1:], "=")); err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }

func (r *session) Run(cmd string) ([]byte, error) {
	// if err := r.SetEnv(); err != nil {
	// 	return nil, err
	// }

	return r.run(func() ([]byte, error) {
		return r.sess.CombinedOutput(cmd)
	})
}

type Output struct {
	output []byte
	err    error
}

func (r *session) run(callback func() ([]byte, error)) ([]byte, error) {
	outputChan := make(chan Output)
	go func() {
		output, err := callback()
		outputChan <- Output{
			output: output,
			err:    err,
		}
	}()

	select {
	case <-r.ctx.Done():
		_ = r.sess.Signal(ssh.SIGINT)

		return nil, r.ctx.Err()
	case result := <-outputChan:
		return result.output, result.err
	}
}
