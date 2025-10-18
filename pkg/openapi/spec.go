package openapi

import (
	"encoding/json"
	"fmt"
)

// Build generates the OpenAPI specification document for the public API.
// The resulting JSON is generated at application startup so the latest
// routes and schemas are always included in the served document.
func Build(apiVersion string) ([]byte, error) {
	doc := map[string]any{
		"openapi": "3.1.0",
		"info": map[string]any{
			"title":       "api.etin.dev",
			"description": "Backend API powering api.etin.dev.",
			"version":     apiVersion,
		},
		"servers": []map[string]any{
			{"url": "https://api.etin.dev"},
			{"url": "http://localhost:4000"},
		},
	}

	components := map[string]any{
		"schemas":         buildSchemas(),
		"securitySchemes": buildSecuritySchemes(),
	}

	doc["components"] = components
	doc["paths"] = buildPaths()

	spec, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("openapi: marshal document: %w", err)
	}

	spec = append(spec, '\n')
	return spec, nil
}

func buildSecuritySchemes() map[string]any {
	return map[string]any{
		"bearerAuth": map[string]any{
			"type":         "http",
			"scheme":       "bearer",
			"bearerFormat": "opaque token",
			"description":  "Bearer token issued by the admin login endpoint.",
		},
	}
}

func buildSchemas() map[string]any {
	stringSchema := func(description string) map[string]any {
		schema := map[string]any{"type": "string"}
		if description != "" {
			schema["description"] = description
		}
		return schema
	}

	dateTimeSchema := func(description string) map[string]any {
		schema := stringSchema(description)
		schema["format"] = "date-time"
		return schema
	}

	int64Schema := func(description string) map[string]any {
		schema := map[string]any{"type": "integer", "format": "int64"}
		if description != "" {
			schema["description"] = description
		}
		return schema
	}

	boolSchema := func(description string) map[string]any {
		schema := map[string]any{"type": "boolean"}
		if description != "" {
			schema["description"] = description
		}
		return schema
	}

	ref := func(name string) map[string]any {
		return map[string]any{"$ref": "#/components/schemas/" + name}
	}

	return map[string]any{
		"AdminLoginRequest": map[string]any{
			"type":     "object",
			"required": []string{"email", "password"},
			"properties": map[string]any{
				"email":    stringSchema("Admin email address."),
				"password": stringSchema("Admin password."),
			},
		},
		"AdminLoginResponse": map[string]any{
			"type":     "object",
			"required": []string{"token", "expiresAt"},
			"properties": map[string]any{
				"token":     stringSchema("Bearer token used to authenticate administrative requests."),
				"expiresAt": dateTimeSchema("Timestamp when the session token expires."),
			},
		},
		"AdminLogoutResponse": map[string]any{
			"type": "object",
			"properties": map[string]any{
				"message": stringSchema("Confirmation message."),
			},
		},
		"HealthcheckResponse": map[string]any{
			"type":     "object",
			"required": []string{"status", "environment", "version"},
			"properties": map[string]any{
				"status":      stringSchema("Service availability status."),
				"environment": stringSchema("Deployment environment for the running service."),
				"version":     stringSchema("Semantic version of the running service."),
			},
		},
		"Company": map[string]any{
			"type":     "object",
			"required": []string{"id", "name", "icon", "description"},
			"properties": map[string]any{
				"id":          int64Schema("Database identifier."),
				"name":        stringSchema("Company display name."),
				"icon":        stringSchema("Icon URL for the company."),
				"description": stringSchema("Markdown description for the company."),
			},
		},
		"CompanyRequest": map[string]any{
			"type": "object",
			"properties": map[string]any{
				"name":        stringSchema("Company display name."),
				"icon":        stringSchema("Icon URL for the company."),
				"description": stringSchema("Markdown description for the company."),
			},
			"required": []string{"name"},
		},
		"UpdateCompanyRequest": map[string]any{
			"type": "object",
			"properties": map[string]any{
				"name":        stringSchema("Company display name."),
				"icon":        stringSchema("Icon URL for the company."),
				"description": stringSchema("Markdown description for the company."),
			},
		},
		"CompanyResponse": map[string]any{
			"type": "object",
			"properties": map[string]any{
				"company": ref("Company"),
			},
		},
		"CompaniesResponse": map[string]any{
			"type": "object",
			"properties": map[string]any{
				"companies": map[string]any{
					"type":  "array",
					"items": ref("Company"),
				},
			},
		},
		"Role": map[string]any{
			"type":     "object",
			"required": []string{"id", "startDate", "endDate", "title", "subtitle", "companyId", "company", "companyIcon", "slug", "description", "skills"},
			"properties": map[string]any{
				"id":          int64Schema("Database identifier."),
				"startDate":   dateTimeSchema("Employment start date."),
				"endDate":     dateTimeSchema("Employment end date. Zero timestamp indicates an ongoing role."),
				"title":       stringSchema("Role title."),
				"subtitle":    stringSchema("Role subtitle."),
				"companyId":   int64Schema("Identifier for the related company."),
				"company":     stringSchema("Resolved company name."),
				"companyIcon": stringSchema("Resolved company icon."),
				"slug":        stringSchema("URL slug for the role."),
				"description": stringSchema("Markdown description of responsibilities."),
				"skills": map[string]any{
					"type":  "array",
					"items": stringSchema("Skill associated with the role."),
				},
			},
		},
		"CreateRoleRequest": map[string]any{
			"type":     "object",
			"required": []string{"startDate", "title", "companyId", "skills"},
			"properties": map[string]any{
				"startDate":   dateTimeSchema("Employment start date."),
				"endDate":     dateTimeSchema("Employment end date. Zero timestamp indicates an ongoing role."),
				"title":       stringSchema("Role title."),
				"subtitle":    stringSchema("Role subtitle."),
				"companyId":   int64Schema("Identifier for the related company."),
				"description": stringSchema("Markdown description of responsibilities."),
				"skills": map[string]any{
					"type":  "array",
					"items": stringSchema("Skill associated with the role."),
				},
			},
		},
		"UpdateRoleRequest": map[string]any{
			"type": "object",
			"properties": map[string]any{
				"startDate":   dateTimeSchema("Employment start date."),
				"endDate":     dateTimeSchema("Employment end date. Zero timestamp indicates an ongoing role."),
				"title":       stringSchema("Role title."),
				"subtitle":    stringSchema("Role subtitle."),
				"companyId":   int64Schema("Identifier for the related company."),
				"description": stringSchema("Markdown description of responsibilities."),
				"skills": map[string]any{
					"type":  "array",
					"items": stringSchema("Skill associated with the role."),
				},
			},
		},
		"RoleResponse": map[string]any{
			"type": "object",
			"properties": map[string]any{
				"role": ref("Role"),
			},
		},
		"RolesResponse": map[string]any{
			"type": "object",
			"properties": map[string]any{
				"roles": map[string]any{
					"type":  "array",
					"items": ref("Role"),
				},
			},
		},
		"Project": map[string]any{
			"type":     "object",
			"required": []string{"id", "startDate", "title", "description"},
			"properties": map[string]any{
				"id":          int64Schema("Database identifier."),
				"startDate":   dateTimeSchema("Project start date."),
				"endDate":     dateTimeSchema("Project end date. Omitted while the project is ongoing."),
				"title":       stringSchema("Project title."),
				"description": stringSchema("Project description."),
				"imageUrl":    stringSchema("Public URL of the project's lead image."),
			},
		},
		"CreateProjectRequest": map[string]any{
			"type":     "object",
			"required": []string{"startDate", "title", "description"},
			"properties": map[string]any{
				"startDate":   dateTimeSchema("Project start date."),
				"endDate":     dateTimeSchema("Project end date. Omitted while the project is ongoing."),
				"title":       stringSchema("Project title."),
				"description": stringSchema("Project description."),
				"imageUrl":    stringSchema("Public URL of the project's lead image."),
			},
		},
		"UpdateProjectRequest": map[string]any{
			"type": "object",
			"properties": map[string]any{
				"startDate":   dateTimeSchema("Project start date."),
				"endDate":     dateTimeSchema("Project end date. Omitted while the project is ongoing."),
				"title":       stringSchema("Project title."),
				"description": stringSchema("Project description."),
				"imageUrl":    stringSchema("Public URL of the project's lead image."),
			},
		},
		"ProjectResponse": map[string]any{
			"type": "object",
			"properties": map[string]any{
				"project": ref("Project"),
			},
		},
		"ProjectsResponse": map[string]any{
			"type": "object",
			"properties": map[string]any{
				"projects": map[string]any{
					"type":  "array",
					"items": ref("Project"),
				},
			},
		},
		"Tag": map[string]any{
			"type":     "object",
			"required": []string{"id", "name", "slug"},
			"properties": map[string]any{
				"id":    int64Schema("Database identifier."),
				"name":  stringSchema("Tag display name."),
				"slug":  stringSchema("URL slug for the tag."),
				"icon":  stringSchema("Optional emoji or icon for the tag."),
				"theme": stringSchema("Optional theme identifier for the tag."),
			},
		},
		"TagsResponse": map[string]any{
			"type": "object",
			"properties": map[string]any{
				"tags": map[string]any{
					"type":  "array",
					"items": ref("Tag"),
				},
			},
		},
		"PublicTag": map[string]any{
			"type":     "object",
			"required": []string{"id", "name", "slug"},
			"properties": map[string]any{
				"id":    int64Schema("Tag identifier."),
				"name":  stringSchema("Tag display name."),
				"slug":  stringSchema("URL-friendly tag slug."),
				"icon":  stringSchema("Optional emoji or icon representing the tag."),
				"theme": stringSchema("Optional theme associated with the tag."),
			},
		},
		"PublicNote": map[string]any{
			"type":     "object",
			"required": []string{"id", "publishedAt", "title", "preview", "body", "isFeatured", "tags"},
			"properties": map[string]any{
				"id":          int64Schema("Note identifier."),
				"publishedAt": stringSchema("ISO 8601 timestamp when the note was published. Empty when unpublished."),
				"title":       stringSchema("Note title."),
				"preview":     stringSchema("Short preview extracted from the note body."),
				"body":        stringSchema("Full note body in Markdown."),
				"isFeatured":  boolSchema("Indicates whether the note is featured."),
				"tags": map[string]any{
					"type":  "array",
					"items": ref("PublicTag"),
				},
			},
		},
		"PublicProject": map[string]any{
			"type":     "object",
			"required": []string{"id", "startDate", "title", "image", "slug", "description", "technologies"},
			"properties": map[string]any{
				"id":        int64Schema("Project identifier."),
				"startDate": dateTimeSchema("Project start date."),
				"endDate": func() map[string]any {
					schema := dateTimeSchema("Project end date. Null while the project is ongoing.")
					schema["nullable"] = true
					return schema
				}(),
				"title": stringSchema("Project title."),
				"image": stringSchema("Lead image URL for the project."),
				"slug":  stringSchema("URL-friendly project slug."),
				"status": func() map[string]any {
					schema := ref("PublicTag")
					schema["nullable"] = true
					schema["description"] = "Status tag assigned to the project."
					return schema
				}(),
				"description": stringSchema("Project description in Markdown."),
				"technologies": map[string]any{
					"type":  "array",
					"items": stringSchema("Technology associated with the project."),
				},
			},
		},
		"PublicRole": map[string]any{
			"type":     "object",
			"required": []string{"roleId", "startDate", "title", "company", "companyIcon", "slug", "description", "skills"},
			"properties": map[string]any{
				"roleId":    int64Schema("Role identifier."),
				"startDate": dateTimeSchema("Role start date."),
				"endDate": func() map[string]any {
					schema := dateTimeSchema("Role end date. Null while the role is ongoing.")
					schema["nullable"] = true
					return schema
				}(),
				"title": stringSchema("Role title."),
				"subtitle": func() map[string]any {
					schema := stringSchema("Optional subtitle providing additional context.")
					schema["nullable"] = true
					return schema
				}(),
				"company":     stringSchema("Company name."),
				"companyIcon": stringSchema("Company icon URL or emoji."),
				"slug":        stringSchema("URL-friendly role slug."),
				"description": stringSchema("Role description in Markdown."),
				"skills": map[string]any{
					"type":  "array",
					"items": stringSchema("Skill associated with the role."),
				},
			},
		},
		"PublicNotesResponse": map[string]any{
			"type": "object",
			"properties": map[string]any{
				"notes": map[string]any{
					"type":  "array",
					"items": ref("PublicNote"),
				},
			},
		},
		"PublicProjectsResponse": map[string]any{
			"type": "object",
			"properties": map[string]any{
				"projects": map[string]any{
					"type":  "array",
					"items": ref("PublicProject"),
				},
			},
		},
		"PublicRolesResponse": map[string]any{
			"type": "object",
			"properties": map[string]any{
				"roles": map[string]any{
					"type":  "array",
					"items": ref("PublicRole"),
				},
			},
		},
		"Note": map[string]any{
			"type":     "object",
			"required": []string{"id", "title", "subtitle", "body"},
			"properties": map[string]any{
				"id":          int64Schema("Database identifier."),
				"publishedAt": dateTimeSchema("Publication timestamp."),
				"title":       stringSchema("Note title."),
				"subtitle":    stringSchema("Note subtitle."),
				"body":        stringSchema("Note body in Markdown."),
			},
		},
		"CreateNoteRequest": map[string]any{
			"type":     "object",
			"required": []string{"title"},
			"properties": map[string]any{
				"title":       stringSchema("Note title."),
				"subtitle":    stringSchema("Note subtitle."),
				"body":        stringSchema("Note body in Markdown."),
				"publishedAt": dateTimeSchema("Publication timestamp."),
			},
		},
		"UpdateNoteRequest": map[string]any{
			"type": "object",
			"properties": map[string]any{
				"title":       stringSchema("Note title."),
				"subtitle":    stringSchema("Note subtitle."),
				"body":        stringSchema("Note body in Markdown."),
				"publishedAt": dateTimeSchema("Publication timestamp."),
			},
		},
		"NoteResponse": map[string]any{
			"type": "object",
			"properties": map[string]any{
				"note": ref("Note"),
			},
		},
		"NotesResponse": map[string]any{
			"type": "object",
			"properties": map[string]any{
				"notes": map[string]any{
					"type":  "array",
					"items": ref("Note"),
				},
			},
		},
		"ItemNote": map[string]any{
			"type":     "object",
			"required": []string{"id", "noteId", "itemId", "itemType"},
			"properties": map[string]any{
				"id":     int64Schema("Database identifier."),
				"noteId": int64Schema("Linked note identifier."),
				"itemId": int64Schema("Linked item identifier."),
				"itemType": map[string]any{
					"type":        "string",
					"description": "Type of item linked to the note.",
					"enum":        []string{"notes", "roles", "projects"},
				},
			},
		},
		"CreateItemNoteRequest": map[string]any{
			"type":     "object",
			"required": []string{"noteId", "itemId", "itemType"},
			"properties": map[string]any{
				"noteId": int64Schema("Linked note identifier."),
				"itemId": int64Schema("Linked item identifier."),
				"itemType": map[string]any{
					"type":        "string",
					"description": "Type of item linked to the note.",
					"enum":        []string{"notes", "roles", "projects"},
				},
			},
		},
		"UpdateItemNoteRequest": map[string]any{
			"type": "object",
			"properties": map[string]any{
				"noteId": int64Schema("Linked note identifier."),
				"itemId": int64Schema("Linked item identifier."),
				"itemType": map[string]any{
					"type":        "string",
					"description": "Type of item linked to the note.",
					"enum":        []string{"notes", "roles", "projects"},
				},
			},
		},
		"ItemNoteResponse": map[string]any{
			"type": "object",
			"properties": map[string]any{
				"itemNote": ref("ItemNote"),
			},
		},
		"ItemNotesResponse": map[string]any{
			"type": "object",
			"properties": map[string]any{
				"itemNotes": map[string]any{
					"type":  "array",
					"items": ref("ItemNote"),
				},
			},
		},
		"TagItem": map[string]any{
			"type":     "object",
			"required": []string{"id", "tagId", "itemId", "itemType"},
			"properties": map[string]any{
				"id":     int64Schema("Database identifier."),
				"tagId":  int64Schema("Linked tag identifier."),
				"itemId": int64Schema("Linked item identifier."),
				"itemType": map[string]any{
					"type":        "string",
					"description": "Type of item linked to the tag.",
					"enum":        []string{"notes", "roles", "projects"},
				},
			},
		},
		"CreateTagItemRequest": map[string]any{
			"type":     "object",
			"required": []string{"tagId", "itemId", "itemType"},
			"properties": map[string]any{
				"tagId":  int64Schema("Linked tag identifier."),
				"itemId": int64Schema("Linked item identifier."),
				"itemType": map[string]any{
					"type":        "string",
					"description": "Type of item linked to the tag.",
					"enum":        []string{"notes", "roles", "projects"},
				},
			},
		},
		"UpdateTagItemRequest": map[string]any{
			"type": "object",
			"properties": map[string]any{
				"tagId":  int64Schema("Linked tag identifier."),
				"itemId": int64Schema("Linked item identifier."),
				"itemType": map[string]any{
					"type":        "string",
					"description": "Type of item linked to the tag.",
					"enum":        []string{"notes", "roles", "projects"},
				},
			},
		},
		"TagItemResponse": map[string]any{
			"type": "object",
			"properties": map[string]any{
				"taggedItem": ref("TagItem"),
			},
		},
		"TagItemsResponse": map[string]any{
			"type": "object",
			"properties": map[string]any{
				"taggedItems": map[string]any{
					"type":  "array",
					"items": ref("TagItem"),
				},
			},
		},
		"Asset": map[string]any{
			"type":     "object",
			"required": []string{"id", "url", "secureUrl", "publicId", "format", "resourceType", "bytes", "width", "height"},
			"properties": map[string]any{
				"id":           int64Schema("Database identifier."),
				"url":          stringSchema("Direct URL for the uploaded asset."),
				"secureUrl":    stringSchema("HTTPS URL for the uploaded asset."),
				"publicId":     stringSchema("Cloudinary public identifier."),
				"format":       stringSchema("File format reported by Cloudinary."),
				"resourceType": stringSchema("Asset resource type such as image or video."),
				"bytes":        int64Schema("File size in bytes."),
				"width":        map[string]any{"type": "integer", "description": "Pixel width when available."},
				"height":       map[string]any{"type": "integer", "description": "Pixel height when available."},
			},
		},
		"AssetUploadRequest": map[string]any{
			"type":     "object",
			"required": []string{"file"},
			"properties": map[string]any{
				"file": map[string]any{
					"type":        "string",
					"format":      "binary",
					"description": "Binary payload of the asset to upload.",
				},
			},
		},
		"AssetUploadResponse": map[string]any{
			"type": "object",
			"properties": map[string]any{
				"asset": ref("Asset"),
			},
		},
	}
}

func buildPaths() map[string]any {
	ref := func(name string) map[string]any {
		return map[string]any{"$ref": "#/components/schemas/" + name}
	}

	jsonResponse := func(description, schema string) map[string]any {
		return map[string]any{
			"description": description,
			"content": map[string]any{
				"application/json": map[string]any{
					"schema": ref(schema),
				},
			},
		}
	}

	noContent := func(description string) map[string]any {
		return map[string]any{"description": description}
	}

	bearerSecurity := []map[string]any{{"bearerAuth": []string{}}}

	intPathParam := func(name, description string) map[string]any {
		return map[string]any{
			"name":        name,
			"in":          "path",
			"required":    true,
			"description": description,
			"schema": map[string]any{
				"type":   "integer",
				"format": "int64",
			},
		}
	}

	itemTypeParam := map[string]any{
		"name":        "itemType",
		"in":          "path",
		"required":    true,
		"description": "Type of item to fetch related records for.",
		"schema": map[string]any{
			"type": "string",
			"enum": []string{"notes", "roles", "projects"},
		},
	}

	itemIdParam := intPathParam("itemId", "Identifier of the item to fetch related records for.")

	paths := map[string]any{
		"/v1/healthcheck": map[string]any{
			"get": map[string]any{
				"operationId": "getHealthcheck",
				"summary":     "Check service health",
				"tags":        []string{"Health"},
				"responses": map[string]any{
					"200": jsonResponse("Service is available.", "HealthcheckResponse"),
				},
			},
		},
		"/swagger": map[string]any{
			"get": map[string]any{
				"operationId": "getSwaggerDocument",
				"summary":     "Retrieve the OpenAPI specification",
				"tags":        []string{"Documentation"},
				"responses": map[string]any{
					"200": map[string]any{
						"description": "OpenAPI document for the API.",
						"content": map[string]any{
							"application/json": map[string]any{
								"schema": map[string]any{"type": "object"},
							},
						},
					},
				},
			},
		},
		"/public/v1/notes": map[string]any{
			"get": map[string]any{
				"operationId": "listPublicNotes",
				"summary":     "List public notes",
				"tags":        []string{"Public Content"},
				"responses": map[string]any{
					"200": jsonResponse("Public notes retrieved.", "PublicNotesResponse"),
					"500": noContent("Server error retrieving public notes."),
				},
			},
		},
		"/public/v1/projects": map[string]any{
			"get": map[string]any{
				"operationId": "listPublicProjects",
				"summary":     "List public projects",
				"tags":        []string{"Public Content"},
				"responses": map[string]any{
					"200": jsonResponse("Public projects retrieved.", "PublicProjectsResponse"),
					"500": noContent("Server error retrieving public projects."),
				},
			},
		},
		"/public/v1/roles": map[string]any{
			"get": map[string]any{
				"operationId": "listPublicRoles",
				"summary":     "List public roles",
				"tags":        []string{"Public Content"},
				"responses": map[string]any{
					"200": jsonResponse("Public roles retrieved.", "PublicRolesResponse"),
					"500": noContent("Server error retrieving public roles."),
				},
			},
		},
		"/v1/admin/login": map[string]any{
			"post": map[string]any{
				"operationId": "adminLogin",
				"summary":     "Create an admin session token",
				"tags":        []string{"Administration"},
				"requestBody": map[string]any{
					"required": true,
					"content": map[string]any{
						"application/json": map[string]any{
							"schema": ref("AdminLoginRequest"),
						},
					},
				},
				"responses": map[string]any{
					"200": jsonResponse("Admin session created.", "AdminLoginResponse"),
					"400": noContent("Invalid credentials payload."),
					"401": noContent("Invalid admin credentials."),
				},
			},
		},
		"/v1/admin/logout": map[string]any{
			"post": map[string]any{
				"operationId": "adminLogout",
				"summary":     "Invalidate the current admin session token",
				"tags":        []string{"Administration"},
				"security":    bearerSecurity,
				"responses": map[string]any{
					"200": jsonResponse("Admin session revoked.", "AdminLogoutResponse"),
					"401": noContent("Missing or invalid bearer token."),
				},
			},
		},
		"/v1/assets": map[string]any{
			"post": map[string]any{
				"operationId": "uploadAsset",
				"summary":     "Upload a new asset",
				"tags":        []string{"Assets"},
				"security":    bearerSecurity,
				"requestBody": map[string]any{
					"required": true,
					"content": map[string]any{
						"multipart/form-data": map[string]any{
							"schema": ref("AssetUploadRequest"),
						},
					},
				},
				"responses": map[string]any{
					"201": jsonResponse("Asset uploaded.", "AssetUploadResponse"),
					"400": noContent("Invalid upload payload or missing file."),
					"403": noContent("Missing or invalid bearer token."),
					"413": noContent("Uploaded file exceeds the maximum allowed size."),
					"500": noContent("Failed to persist asset metadata."),
					"502": noContent("Failed to upload asset to storage provider."),
				},
			},
		},
		"/v1/roles": map[string]any{
			"get": map[string]any{
				"operationId": "listRoles",
				"summary":     "List roles",
				"tags":        []string{"Roles"},
				"responses": map[string]any{
					"200": jsonResponse("Roles retrieved.", "RolesResponse"),
					"500": noContent("Server error retrieving roles."),
				},
			},
			"post": map[string]any{
				"operationId": "createRole",
				"summary":     "Create a role",
				"tags":        []string{"Roles"},
				"security":    bearerSecurity,
				"requestBody": map[string]any{
					"required": true,
					"content": map[string]any{
						"application/json": map[string]any{
							"schema": ref("CreateRoleRequest"),
						},
					},
				},
				"responses": map[string]any{
					"201": jsonResponse("Role created.", "RoleResponse"),
					"400": noContent("Invalid payload."),
					"403": noContent("Missing or invalid bearer token."),
				},
			},
		},
		"/v1/roles/{roleId}": map[string]any{
			"get": map[string]any{
				"operationId": "getRole",
				"summary":     "Retrieve a role",
				"tags":        []string{"Roles"},
				"parameters":  []map[string]any{intPathParam("roleId", "Identifier of the role.")},
				"responses": map[string]any{
					"200": jsonResponse("Role retrieved.", "RoleResponse"),
					"400": noContent("Invalid role identifier."),
					"404": noContent("Role not found."),
				},
			},
			"put": map[string]any{
				"operationId": "updateRole",
				"summary":     "Update a role",
				"tags":        []string{"Roles"},
				"security":    bearerSecurity,
				"parameters":  []map[string]any{intPathParam("roleId", "Identifier of the role.")},
				"requestBody": map[string]any{
					"required": true,
					"content": map[string]any{
						"application/json": map[string]any{
							"schema": ref("UpdateRoleRequest"),
						},
					},
				},
				"responses": map[string]any{
					"200": jsonResponse("Role updated.", "RoleResponse"),
					"400": noContent("Invalid payload."),
					"403": noContent("Missing or invalid bearer token."),
					"404": noContent("Role not found."),
					"500": noContent("Server error updating role."),
				},
			},
			"delete": map[string]any{
				"operationId": "deleteRole",
				"summary":     "Delete a role",
				"tags":        []string{"Roles"},
				"security":    bearerSecurity,
				"parameters":  []map[string]any{intPathParam("roleId", "Identifier of the role.")},
				"responses": map[string]any{
					"204": noContent("Role deleted."),
					"400": noContent("Invalid role identifier."),
					"403": noContent("Missing or invalid bearer token."),
					"404": noContent("Role not found."),
				},
			},
		},
		"/v1/companies": map[string]any{
			"get": map[string]any{
				"operationId": "listCompanies",
				"summary":     "List companies",
				"tags":        []string{"Companies"},
				"responses": map[string]any{
					"200": jsonResponse("Companies retrieved.", "CompaniesResponse"),
					"500": noContent("Server error retrieving companies."),
				},
			},
			"post": map[string]any{
				"operationId": "createCompany",
				"summary":     "Create a company",
				"tags":        []string{"Companies"},
				"security":    bearerSecurity,
				"requestBody": map[string]any{
					"required": true,
					"content": map[string]any{
						"application/json": map[string]any{
							"schema": ref("UpdateCompanyRequest"),
						},
					},
				},
				"responses": map[string]any{
					"200": jsonResponse("Company created.", "CompanyResponse"),
					"400": noContent("Invalid payload."),
					"403": noContent("Missing or invalid bearer token."),
				},
			},
		},
		"/v1/companies/{companyId}": map[string]any{
			"get": map[string]any{
				"operationId": "getCompany",
				"summary":     "Retrieve a company",
				"tags":        []string{"Companies"},
				"parameters":  []map[string]any{intPathParam("companyId", "Identifier of the company.")},
				"responses": map[string]any{
					"200": jsonResponse("Company retrieved.", "CompanyResponse"),
					"400": noContent("Invalid company identifier."),
					"500": noContent("Server error retrieving company."),
				},
			},
			"put": map[string]any{
				"operationId": "updateCompany",
				"summary":     "Update a company",
				"tags":        []string{"Companies"},
				"security":    bearerSecurity,
				"parameters":  []map[string]any{intPathParam("companyId", "Identifier of the company.")},
				"requestBody": map[string]any{
					"required": true,
					"content": map[string]any{
						"application/json": map[string]any{
							"schema": ref("UpdateCompanyRequest"),
						},
					},
				},
				"responses": map[string]any{
					"202": jsonResponse("Company updated.", "CompanyResponse"),
					"400": noContent("Invalid payload."),
					"403": noContent("Missing or invalid bearer token."),
					"500": noContent("Server error updating company."),
				},
			},
			"delete": map[string]any{
				"operationId": "deleteCompany",
				"summary":     "Delete a company",
				"tags":        []string{"Companies"},
				"security":    bearerSecurity,
				"parameters":  []map[string]any{intPathParam("companyId", "Identifier of the company.")},
				"responses": map[string]any{
					"204": noContent("Company deleted."),
					"400": noContent("Invalid company identifier."),
					"403": noContent("Missing or invalid bearer token."),
					"500": noContent("Server error deleting company."),
				},
			},
		},
		"/v1/projects": map[string]any{
			"get": map[string]any{
				"operationId": "listProjects",
				"summary":     "List projects",
				"tags":        []string{"Projects"},
				"responses": map[string]any{
					"200": jsonResponse("Projects retrieved.", "ProjectsResponse"),
					"500": noContent("Server error retrieving projects."),
				},
			},
			"post": map[string]any{
				"operationId": "createProject",
				"summary":     "Create a project",
				"tags":        []string{"Projects"},
				"security":    bearerSecurity,
				"requestBody": map[string]any{
					"required": true,
					"content": map[string]any{
						"application/json": map[string]any{
							"schema": ref("CreateProjectRequest"),
						},
					},
				},
				"responses": map[string]any{
					"201": jsonResponse("Project created.", "ProjectResponse"),
					"400": noContent("Invalid payload."),
					"403": noContent("Missing or invalid bearer token."),
				},
			},
		},
		"/v1/projects/{projectId}": map[string]any{
			"get": map[string]any{
				"operationId": "getProject",
				"summary":     "Retrieve a project",
				"tags":        []string{"Projects"},
				"parameters":  []map[string]any{intPathParam("projectId", "Identifier of the project.")},
				"responses": map[string]any{
					"200": jsonResponse("Project retrieved.", "ProjectResponse"),
					"400": noContent("Invalid project identifier."),
					"404": noContent("Project not found."),
				},
			},
			"put": map[string]any{
				"operationId": "updateProject",
				"summary":     "Update a project",
				"tags":        []string{"Projects"},
				"security":    bearerSecurity,
				"parameters":  []map[string]any{intPathParam("projectId", "Identifier of the project.")},
				"requestBody": map[string]any{
					"required": true,
					"content": map[string]any{
						"application/json": map[string]any{
							"schema": ref("UpdateProjectRequest"),
						},
					},
				},
				"responses": map[string]any{
					"200": jsonResponse("Project updated.", "ProjectResponse"),
					"400": noContent("Invalid payload."),
					"403": noContent("Missing or invalid bearer token."),
					"404": noContent("Project not found."),
					"500": noContent("Server error updating project."),
				},
			},
			"delete": map[string]any{
				"operationId": "deleteProject",
				"summary":     "Delete a project",
				"tags":        []string{"Projects"},
				"security":    bearerSecurity,
				"parameters":  []map[string]any{intPathParam("projectId", "Identifier of the project.")},
				"responses": map[string]any{
					"204": noContent("Project deleted."),
					"400": noContent("Invalid project identifier."),
					"403": noContent("Missing or invalid bearer token."),
					"404": noContent("Project not found."),
				},
			},
		},
		"/v1/notes": map[string]any{
			"get": map[string]any{
				"operationId": "listNotes",
				"summary":     "List notes",
				"tags":        []string{"Notes"},
				"responses": map[string]any{
					"200": jsonResponse("Notes retrieved.", "NotesResponse"),
					"500": noContent("Server error retrieving notes."),
				},
			},
			"post": map[string]any{
				"operationId": "createNote",
				"summary":     "Create a note",
				"tags":        []string{"Notes"},
				"security":    bearerSecurity,
				"requestBody": map[string]any{
					"required": true,
					"content": map[string]any{
						"application/json": map[string]any{
							"schema": ref("CreateNoteRequest"),
						},
					},
				},
				"responses": map[string]any{
					"201": jsonResponse("Note created.", "NoteResponse"),
					"400": noContent("Invalid payload."),
					"403": noContent("Missing or invalid bearer token."),
				},
			},
		},
		"/v1/notes/{noteId}": map[string]any{
			"get": map[string]any{
				"operationId": "getNote",
				"summary":     "Retrieve a note",
				"tags":        []string{"Notes"},
				"parameters":  []map[string]any{intPathParam("noteId", "Identifier of the note.")},
				"responses": map[string]any{
					"200": jsonResponse("Note retrieved.", "NoteResponse"),
					"400": noContent("Invalid note identifier."),
					"404": noContent("Note not found."),
				},
			},
			"put": map[string]any{
				"operationId": "updateNote",
				"summary":     "Update a note",
				"tags":        []string{"Notes"},
				"security":    bearerSecurity,
				"parameters":  []map[string]any{intPathParam("noteId", "Identifier of the note.")},
				"requestBody": map[string]any{
					"required": true,
					"content": map[string]any{
						"application/json": map[string]any{
							"schema": ref("UpdateNoteRequest"),
						},
					},
				},
				"responses": map[string]any{
					"200": jsonResponse("Note updated.", "NoteResponse"),
					"400": noContent("Invalid payload."),
					"403": noContent("Missing or invalid bearer token."),
					"404": noContent("Note not found."),
					"500": noContent("Server error updating note."),
				},
			},
			"delete": map[string]any{
				"operationId": "deleteNote",
				"summary":     "Delete a note",
				"tags":        []string{"Notes"},
				"security":    bearerSecurity,
				"parameters":  []map[string]any{intPathParam("noteId", "Identifier of the note.")},
				"responses": map[string]any{
					"204": noContent("Note deleted."),
					"400": noContent("Invalid note identifier."),
					"403": noContent("Missing or invalid bearer token."),
					"404": noContent("Note not found."),
				},
			},
		},
		"/v1/item-notes": map[string]any{
			"get": map[string]any{
				"operationId": "listItemNotes",
				"summary":     "List item-note links",
				"tags":        []string{"Item Notes"},
				"responses": map[string]any{
					"200": jsonResponse("Item note associations retrieved.", "ItemNotesResponse"),
					"500": noContent("Server error retrieving item note associations."),
				},
			},
			"post": map[string]any{
				"operationId": "createItemNote",
				"summary":     "Create an item-note link",
				"tags":        []string{"Item Notes"},
				"security":    bearerSecurity,
				"requestBody": map[string]any{
					"required": true,
					"content": map[string]any{
						"application/json": map[string]any{
							"schema": ref("CreateItemNoteRequest"),
						},
					},
				},
				"responses": map[string]any{
					"201": jsonResponse("Item note association created.", "ItemNoteResponse"),
					"400": noContent("Invalid payload."),
					"403": noContent("Missing or invalid bearer token."),
				},
			},
		},
		"/v1/item-notes/{itemNoteId}": map[string]any{
			"get": map[string]any{
				"operationId": "getItemNote",
				"summary":     "Retrieve an item-note link",
				"tags":        []string{"Item Notes"},
				"parameters":  []map[string]any{intPathParam("itemNoteId", "Identifier of the item-note link.")},
				"responses": map[string]any{
					"200": jsonResponse("Item note association retrieved.", "ItemNoteResponse"),
					"400": noContent("Invalid item-note identifier."),
					"404": noContent("Item note association not found."),
				},
			},
			"put": map[string]any{
				"operationId": "updateItemNote",
				"summary":     "Update an item-note link",
				"tags":        []string{"Item Notes"},
				"security":    bearerSecurity,
				"parameters":  []map[string]any{intPathParam("itemNoteId", "Identifier of the item-note link.")},
				"requestBody": map[string]any{
					"required": true,
					"content": map[string]any{
						"application/json": map[string]any{
							"schema": ref("UpdateItemNoteRequest"),
						},
					},
				},
				"responses": map[string]any{
					"200": jsonResponse("Item note association updated.", "ItemNoteResponse"),
					"400": noContent("Invalid payload."),
					"403": noContent("Missing or invalid bearer token."),
					"404": noContent("Item note association not found."),
					"500": noContent("Server error updating item note association."),
				},
			},
			"delete": map[string]any{
				"operationId": "deleteItemNote",
				"summary":     "Delete an item-note link",
				"tags":        []string{"Item Notes"},
				"security":    bearerSecurity,
				"parameters":  []map[string]any{intPathParam("itemNoteId", "Identifier of the item-note link.")},
				"responses": map[string]any{
					"204": noContent("Item note association deleted."),
					"400": noContent("Invalid item-note identifier."),
					"403": noContent("Missing or invalid bearer token."),
					"404": noContent("Item note association not found."),
				},
			},
		},
		"/v1/item-notes/items/{itemType}/{itemId}": map[string]any{
			"get": map[string]any{
				"operationId": "listNotesForItem",
				"summary":     "List notes associated with an item",
				"tags":        []string{"Item Notes"},
				"parameters":  []map[string]any{itemTypeParam, itemIdParam},
				"responses": map[string]any{
					"200": jsonResponse("Notes retrieved.", "NotesResponse"),
					"400": noContent("Invalid item type or identifier."),
					"500": noContent("Server error retrieving notes for the item."),
				},
			},
		},
		"/v1/tagged-items": map[string]any{
			"get": map[string]any{
				"operationId": "listTagItems",
				"summary":     "List tag associations",
				"tags":        []string{"Tag Items"},
				"responses": map[string]any{
					"200": jsonResponse("Tag associations retrieved.", "TagItemsResponse"),
					"500": noContent("Server error retrieving tag associations."),
				},
			},
			"post": map[string]any{
				"operationId": "createTagItem",
				"summary":     "Create a tag association",
				"tags":        []string{"Tag Items"},
				"security":    bearerSecurity,
				"requestBody": map[string]any{
					"required": true,
					"content": map[string]any{
						"application/json": map[string]any{
							"schema": ref("CreateTagItemRequest"),
						},
					},
				},
				"responses": map[string]any{
					"201": jsonResponse("Tag association created.", "TagItemResponse"),
					"400": noContent("Invalid payload."),
					"403": noContent("Missing or invalid bearer token."),
				},
			},
		},
		"/v1/tagged-items/{taggedItemId}": map[string]any{
			"get": map[string]any{
				"operationId": "getTagItem",
				"summary":     "Retrieve a tag association",
				"tags":        []string{"Tag Items"},
				"parameters":  []map[string]any{intPathParam("taggedItemId", "Identifier of the tag association.")},
				"responses": map[string]any{
					"200": jsonResponse("Tag association retrieved.", "TagItemResponse"),
					"400": noContent("Invalid tag association identifier."),
					"404": noContent("Tag association not found."),
				},
			},
			"put": map[string]any{
				"operationId": "updateTagItem",
				"summary":     "Update a tag association",
				"tags":        []string{"Tag Items"},
				"security":    bearerSecurity,
				"parameters":  []map[string]any{intPathParam("taggedItemId", "Identifier of the tag association.")},
				"requestBody": map[string]any{
					"required": true,
					"content": map[string]any{
						"application/json": map[string]any{
							"schema": ref("UpdateTagItemRequest"),
						},
					},
				},
				"responses": map[string]any{
					"200": jsonResponse("Tag association updated.", "TagItemResponse"),
					"400": noContent("Invalid payload."),
					"403": noContent("Missing or invalid bearer token."),
					"404": noContent("Tag association not found."),
					"500": noContent("Server error updating tag association."),
				},
			},
			"delete": map[string]any{
				"operationId": "deleteTagItem",
				"summary":     "Delete a tag association",
				"tags":        []string{"Tag Items"},
				"security":    bearerSecurity,
				"parameters":  []map[string]any{intPathParam("taggedItemId", "Identifier of the tag association.")},
				"responses": map[string]any{
					"204": noContent("Tag association deleted."),
					"400": noContent("Invalid tag association identifier."),
					"403": noContent("Missing or invalid bearer token."),
					"404": noContent("Tag association not found."),
				},
			},
		},
		"/v1/tagged-items/items/{itemType}/{itemId}": map[string]any{
			"get": map[string]any{
				"operationId": "listTagsForItem",
				"summary":     "List tags associated with an item",
				"tags":        []string{"Tag Items"},
				"parameters":  []map[string]any{itemTypeParam, itemIdParam},
				"responses": map[string]any{
					"200": jsonResponse("Tags retrieved.", "TagsResponse"),
					"400": noContent("Invalid item type or identifier."),
					"500": noContent("Server error retrieving tags for the item."),
				},
			},
		},
	}

	return paths
}
