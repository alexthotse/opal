cat << 'GO' >> /workspace/peregrine/adapters/backend_client.go

func (b *BackendClient) StartUltraplan(ctx context.Context, id string) (string, error) {
        req := connect.NewRequest(&falconv1.UltraplanRequest{Id: id})
        res, err := b.client.StartUltraplan(ctx, req)
        if err != nil {
                return "", err
        }
        return res.Msg.Result, nil
}

func (b *BackendClient) QuickSearch(ctx context.Context, id, query string) (string, error) {
        req := connect.NewRequest(&falconv1.SearchRequest{Id: id, Query: query})
        res, err := b.client.QuickSearch(ctx, req)
        if err != nil {
                return "", err
        }
        return res.Msg.Result, nil
}

func (b *BackendClient) GetTokenBudget(ctx context.Context, id string) (string, error) {
        req := connect.NewRequest(&falconv1.BudgetRequest{Id: id})
        res, err := b.client.GetTokenBudget(ctx, req)
        if err != nil {
                return "", err
        }
        return res.Msg.Result, nil
}

func (b *BackendClient) GetTeamMem(ctx context.Context, id string) (string, error) {
        req := connect.NewRequest(&falconv1.TeamMemRequest{Id: id})
        res, err := b.client.GetTeamMem(ctx, req)
        if err != nil {
                return "", err
        }
        return res.Msg.Result, nil
}

func (b *BackendClient) StartBridgeMode(ctx context.Context, id string) (string, error) {
        req := connect.NewRequest(&falconv1.BridgeModeRequest{Id: id})
        res, err := b.client.StartBridgeMode(ctx, req)
        if err != nil {
                return "", err
        }
        return res.Msg.Result, nil
}

func (b *BackendClient) SetAgentTrigger(ctx context.Context, id string) (string, error) {
        req := connect.NewRequest(&falconv1.TriggerRequest{Id: id})
        res, err := b.client.SetAgentTrigger(ctx, req)
        if err != nil {
                return "", err
        }
        return res.Msg.Result, nil
}

func (b *BackendClient) SetAgentTriggerRemote(ctx context.Context, id string) (string, error) {
        req := connect.NewRequest(&falconv1.TriggerRemoteRequest{Id: id})
        res, err := b.client.SetAgentTriggerRemote(ctx, req)
        if err != nil {
                return "", err
        }
        return res.Msg.Result, nil
}

func (b *BackendClient) ExtractMemories(ctx context.Context, id string) (string, error) {
        req := connect.NewRequest(&falconv1.ExtractMemoriesRequest{Id: id})
        res, err := b.client.ExtractMemories(ctx, req)
        if err != nil {
                return "", err
        }
        return res.Msg.Result, nil
}

func (b *BackendClient) GetCompactionReminders(ctx context.Context, id string) (string, error) {
        req := connect.NewRequest(&falconv1.CompactionRequest{Id: id})
        res, err := b.client.GetCompactionReminders(ctx, req)
        if err != nil {
                return "", err
        }
        return res.Msg.Result, nil
}

func (b *BackendClient) CachedMicrocompact(ctx context.Context, id string) (string, error) {
        req := connect.NewRequest(&falconv1.MicrocompactRequest{Id: id})
        res, err := b.client.CachedMicrocompact(ctx, req)
        if err != nil {
                return "", err
        }
        return res.Msg.Result, nil
}

func (b *BackendClient) GetShotStats(ctx context.Context, id string) (string, error) {
        req := connect.NewRequest(&falconv1.ShotStatsRequest{Id: id})
        res, err := b.client.GetShotStats(ctx, req)
        if err != nil {
                return "", err
        }
        return res.Msg.Result, nil
}

func (b *BackendClient) StartVerification(ctx context.Context, id string) (string, error) {
        req := connect.NewRequest(&falconv1.VerificationRequest{Id: id})
        res, err := b.client.StartVerification(ctx, req)
        if err != nil {
                return "", err
        }
        return res.Msg.Result, nil
}

func (b *BackendClient) ClassifyBash(ctx context.Context, id string) (string, error) {
        req := connect.NewRequest(&falconv1.BashClassifyRequest{Id: id})
        res, err := b.client.ClassifyBash(ctx, req)
        if err != nil {
                return "", err
        }
        return res.Msg.Result, nil
}

func (b *BackendClient) GetVoiceMode(ctx context.Context, id string) (string, error) {
        req := connect.NewRequest(&falconv1.VoiceModeRequest{Id: id})
        res, err := b.client.GetVoiceMode(ctx, req)
        if err != nil {
                return "", err
        }
        return res.Msg.Result, nil
}

func (b *BackendClient) GetHistoryPicker(ctx context.Context, id string) (string, error) {
        req := connect.NewRequest(&falconv1.HistoryPickerRequest{Id: id})
        res, err := b.client.GetHistoryPicker(ctx, req)
        if err != nil {
                return "", err
        }
        return res.Msg.Result, nil
}

func (b *BackendClient) MessageActions(ctx context.Context, id string) (string, error) {
        req := connect.NewRequest(&falconv1.MessageActionsRequest{Id: id})
        res, err := b.client.MessageActions(ctx, req)
        if err != nil {
                return "", err
        }
        return res.Msg.Result, nil
}
GO
