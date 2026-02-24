// CUE Schema for gion configuration
// This schema defines validation rules for gion.yaml and repository specifications

// ============== Provider ==============
#Provider: "github" | "gitlab" | "bitbucket" | "custom"

// ============== EndPoint ==============
#EndPoint: {
    host: string & !=""
    port?: string & ~"^[1-9][0-9]{0,4}$"
    basePath?: string & strings.HasPrefix(_, "/")
}

// ============== Registry ==============
#Registry: {
    provider: #Provider
    group: string & !=""
    subGroups?: [...string]
}

// ============== Repository ==============
#Repository: {
    repo: string & !="" & !~"\\.git$"
}

// ============== Spec (repospec) ==============
#Spec: {
    endPoint: #EndPoint
    registry: #Registry
    repository: #Repository
    repoKey: string & !=""
    isSSH: bool
}

// ============== gion.yaml Schema ==============

// Workspace Mode (from INVENTORY.md)
#WorkspaceMode: "preset" | "repo" | "review" | "issue" | "resume" | "add"

// Preset Repository Entry
#PresetRepoEntry: {
    repo: string & !=""
    provider?: #Provider
    base_path?: string & strings.HasPrefix(_, "/")
}

// Preset Definition
#Preset: {
    repos: [...#PresetRepoEntry]
}

// Workspace Repository Entry
#WorkspaceRepoEntry: {
    alias: string & !=""
    repo_key: string & !=""
    branch: string & !=""
    base_ref?: string & strings.HasPrefix(_, "origin/")
}

// Workspace Definition
#Workspace: {
    description?: string
    mode: #WorkspaceMode
    source_url?: string
    preset_name?: string
    repos: [...#WorkspaceRepoEntry]
}

// Root Configuration (matches INVENTORY.md)
#Config: {
    version: int & >= 1
    presets?: [string]: #Preset
    workspaces: [string]: #Workspace
}
