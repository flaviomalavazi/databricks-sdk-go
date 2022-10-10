package internal

import (
	"context"
	"testing"

	"github.com/databricks/databricks-sdk-go/service/gitcredentials"
	"github.com/databricks/databricks-sdk-go/workspaces"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAccGitCredentials(t *testing.T) {
	env := GetEnvOrSkipTest(t, "CLOUD_ENV")
	t.Log(env)
	ctx := context.Background()
	wsc := workspaces.New()

	list, err := wsc.GitCredentials.List(ctx)
	require.NoError(t, err)
	for _, v := range list.Credentials {
		err = wsc.GitCredentials.DeleteByCredentialId(ctx, v.CredentialId)
		require.NoError(t, err)
	}

	cr, err := wsc.GitCredentials.Create(ctx, gitcredentials.CreateCredentials{
		GitProvider:         "gitHub",
		GitUsername:         "test",
		PersonalAccessToken: "test",
	})
	require.NoError(t, err)
	t.Cleanup(func() {
		err = wsc.GitCredentials.DeleteByCredentialId(ctx, cr.CredentialId)
		require.NoError(t, err)
	})

	err = wsc.GitCredentials.Update(ctx, gitcredentials.UpdateCredentials{
		CredentialId:        cr.CredentialId,
		GitProvider:         "gitHub",
		GitUsername:         RandomEmail(),
		PersonalAccessToken: RandomName(),
	})
	require.NoError(t, err)

	load, err := wsc.GitCredentials.GetByCredentialId(ctx, cr.CredentialId)
	require.NoError(t, err)

	assert.NotEqual(t, cr.GitUsername, load.GitUsername)
}