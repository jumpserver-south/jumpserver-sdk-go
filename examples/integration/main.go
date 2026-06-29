// Integration test: full CRUD lifecycle for every JumpServer v4 service.
//
// Environment variables required:
//
//	JUMPSERVER_URL        — base URL
//	JUMPSERVER_KEY_ID     — access key ID
//	JUMPSERVER_SECRET_ID  — access key secret
//
// Run:
//
//	go run ./examples/integration
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	jumpserver "github.com/jumpserver-south/jumpserver-sdk-go"
	"github.com/jumpserver-south/jumpserver-sdk-go/assets"
	"github.com/jumpserver-south/jumpserver-sdk-go/model"
)

var (
	passed  int
	failed  int
	skipped int
)

func ok(name string) {
	fmt.Printf("  ✓ %s\n", name)
	passed++
}

func fail(name string, err error) {
	fmt.Printf("  ✗ %-50s %s\n", name, shortErr(err))
	failed++
}

func skip(name, reason string) {
	fmt.Printf("  ○ %-50s [SKIP: %s]\n", name, reason)
	skipped++
}

func shortErr(err error) string {
	if err == nil {
		return ""
	}
	s := err.Error()
	if len(s) > 160 {
		return s[:160] + "..."
	}
	return s
}

func section(title string) {
	fmt.Printf("\n=== %s ===\n", title)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func main() {
	url := os.Getenv("JUMPSERVER_URL")
	keyID := os.Getenv("JUMPSERVER_KEY_ID")
	secretID := os.Getenv("JUMPSERVER_SECRET_ID")
	if url == "" || keyID == "" || secretID == "" {
		log.Fatal("JUMPSERVER_URL, JUMPSERVER_KEY_ID, JUMPSERVER_SECRET_ID required")
	}

	client := jumpserver.NewClient(
		jumpserver.WithBaseURL(url),
		jumpserver.WithAccessKeyAuth(keyID, secretID),
	)
	ctx := context.Background()
	ts := fmt.Sprintf("%d", time.Now().UnixNano()%100000)

	// Get a valid org ID for scoped operations.
	orgs, _, err := client.Organizations.List(ctx, &jumpserver.ListOptions{Limit: 10})
	var orgID string
	if err == nil && len(orgs) > 0 {
		for _, o := range orgs {
			if !o.IsRoot {
				orgID = o.ID
				break
			}
		}
		if orgID == "" {
			orgID = orgs[0].ID
		}
	}
	if orgID == "" {
		log.Fatal("No organizations found")
	}
	fmt.Printf("Using org: %s\n", orgID)
	scoped := client.WithOrgScope(orgID)

	var (
		createdNodeID      string
		createdZoneID      string
		createdLabelID     string
		createdHostID      string
		createdAccountID   string
		createdTemplateID  string
		createdPermID      string
		createdCmdFilterID string
		createdCmdGroupID  string
		createdCategoryIDs = make(map[string]string) // category name -> asset ID
	)

	// ============================================================
	section("Settings")
	{
		pub, _, err := client.Settings.Public(ctx)
		if err != nil {
			fail("Settings.Public", err)
		} else {
			ok(fmt.Sprintf("Settings.Public (watermark=%v)", pub.EnableWatermark))
		}
		settings, _, err := client.Settings.List(ctx, nil)
		if err != nil {
			fail("Settings.List", err)
		} else {
			ok(fmt.Sprintf("Settings.List (%d keys)", len(settings)))
		}
	}

	// ============================================================
	section("Users")
	{
		profile, _, err := client.Users.Profile(ctx)
		if err != nil {
			fail("Users.Profile", err)
		} else {
			ok(fmt.Sprintf("Users.Profile (%s)", profile.Username))
		}
		users, _, err := client.Users.List(ctx, nil, &jumpserver.ListOptions{Limit: 15})
		if err != nil {
			fail("Users.List", err)
		} else {
			ok(fmt.Sprintf("Users.List (%d)", len(users)))
		}
		if len(users) > 0 {
			u, _, err := client.Users.Get(ctx, users[0].ID)
			if err != nil {
				fail("Users.Get", err)
			} else {
				ok(fmt.Sprintf("Users.Get (%s)", u.Username))
			}
		}
	}

	// ============================================================
	section("UserGroups CRUD")
	{
		g, _, err := scoped.UserGroups.Create(ctx, &model.GroupRequest{Name: "ug-" + ts})
		if err != nil {
			fail("UserGroups.Create", err)
		} else {
			ok(fmt.Sprintf("UserGroups.Create (id=%s)", g.ID))
			u, _, err := scoped.UserGroups.Update(ctx, g.ID, &model.GroupRequest{Name: "ug-upd-" + ts})
			if err != nil {
				fail("UserGroups.Update", err)
			} else {
				ok(fmt.Sprintf("UserGroups.Update (%s)", u.Name))
			}
			got, _, err := scoped.UserGroups.Get(ctx, g.ID)
			if err != nil {
				fail("UserGroups.Get", err)
			} else {
				ok(fmt.Sprintf("UserGroups.Get (%s)", got.Name))
			}
			list, _, err := scoped.UserGroups.List(ctx, &jumpserver.ListOptions{Limit: 15})
			if err != nil {
				fail("UserGroups.List", err)
			} else {
				ok(fmt.Sprintf("UserGroups.List (%d)", len(list)))
			}
			lu, _, err := scoped.UserGroups.ListUsers(ctx, g.ID, &jumpserver.ListOptions{Limit: 15})
			if err != nil {
				skip("UserGroups.ListUsers", err.Error())
			} else {
				ok(fmt.Sprintf("UserGroups.ListUsers (%d)", len(lu)))
			}
			_, err = scoped.UserGroups.Delete(ctx, g.ID)
			if err != nil {
				fail("UserGroups.Delete", err)
			} else {
				ok("UserGroups.Delete")
			}
		}
	}

	// ============================================================
	section("Roles")
	{
		orgRoles, _, err := client.Roles.List(ctx, "org", &jumpserver.ListOptions{Limit: 15})
		if err != nil {
			fail("Roles.List(org)", err)
		} else {
			ok(fmt.Sprintf("Roles.List(org) (%d)", len(orgRoles)))
		}
		sysRoles, _, err := client.Roles.List(ctx, "system", &jumpserver.ListOptions{Limit: 15})
		if err != nil {
			fail("Roles.List(system)", err)
		} else {
			ok(fmt.Sprintf("Roles.List(system) (%d)", len(sysRoles)))
		}
		if len(orgRoles) > 0 {
			r, _, err := client.Roles.Get(ctx, "org", orgRoles[0].ID)
			if err != nil {
				fail("Roles.Get(org)", err)
			} else {
				ok(fmt.Sprintf("Roles.Get(org) (%s)", r.Name))
			}
		}
	}

	// ============================================================
	section("Organizations CRUD")
	{
		o, _, err := client.Organizations.Create(ctx, &model.OrganizationRequest{Name: "org-" + ts})
		if err != nil {
			fail("Organizations.Create", err)
		} else {
			ok(fmt.Sprintf("Organizations.Create (id=%s)", o.ID))
			u, _, err := client.Organizations.Update(ctx, o.ID, &model.OrganizationRequest{Name: "org-upd-" + ts})
			if err != nil {
				fail("Organizations.Update", err)
			} else {
				ok(fmt.Sprintf("Organizations.Update (%s)", u.Name))
			}
			got, _, err := client.Organizations.Get(ctx, o.ID)
			if err != nil {
				fail("Organizations.Get", err)
			} else {
				ok(fmt.Sprintf("Organizations.Get (%s)", got.Name))
			}
			list, _, err := client.Organizations.List(ctx, &jumpserver.ListOptions{Limit: 15})
			if err != nil {
				fail("Organizations.List", err)
			} else {
				ok(fmt.Sprintf("Organizations.List (%d)", len(list)))
			}
			_, err = client.Organizations.Delete(ctx, o.ID)
			if err != nil {
				fail("Organizations.Delete", err)
			} else {
				ok("Organizations.Delete")
			}
		}
	}

	// ============================================================
	section("Platforms")
	{
		platforms, _, err := client.Platforms.List(ctx, &jumpserver.ListOptions{Limit: 15})
		if err != nil {
			fail("Platforms.List", err)
		} else {
			ok(fmt.Sprintf("Platforms.List (%d)", len(platforms)))
		}
		if len(platforms) > 0 {
			p, _, err := client.Platforms.Get(ctx, platforms[0].ID)
			if err != nil {
				fail("Platforms.Get", err)
			} else {
				ok(fmt.Sprintf("Platforms.Get (%s)", p.Name))
			}
		}
	}

	// ============================================================
	section("Nodes CRUD")
	{
		n, _, err := scoped.Nodes.Create(ctx, &model.NodeRequest{Value: "node-" + ts})
		if err != nil {
			fail("Nodes.Create", err)
		} else {
			createdNodeID = n.ID
			ok(fmt.Sprintf("Nodes.Create (id=%s)", n.ID))
			u, _, err := scoped.Nodes.Update(ctx, n.ID, &model.NodeRequest{Value: "node-upd-" + ts})
			if err != nil {
				fail("Nodes.Update", err)
			} else {
				ok(fmt.Sprintf("Nodes.Update (%s)", u.Value))
			}
			got, _, err := scoped.Nodes.Get(ctx, n.ID)
			if err != nil {
				fail("Nodes.Get", err)
			} else {
				ok(fmt.Sprintf("Nodes.Get (%s)", got.Value))
			}
			list, _, err := scoped.Nodes.List(ctx, &jumpserver.ListOptions{Limit: 15})
			if err != nil {
				fail("Nodes.List", err)
			} else {
				ok(fmt.Sprintf("Nodes.List (%d)", len(list)))
			}
		}
	}

	// ============================================================
	section("Zones CRUD")
	{
		z, _, err := scoped.Zones.Create(ctx, &model.ZoneRequest{Name: "zone-" + ts})
		if err != nil {
			fail("Zones.Create", err)
		} else {
			createdZoneID = z.ID
			ok(fmt.Sprintf("Zones.Create (id=%s)", z.ID))
			u, _, err := scoped.Zones.Update(ctx, z.ID, &model.ZoneRequest{Name: "zone-upd-" + ts})
			if err != nil {
				fail("Zones.Update", err)
			} else {
				ok(fmt.Sprintf("Zones.Update (%s)", u.Name))
			}
			got, _, err := scoped.Zones.Get(ctx, z.ID)
			if err != nil {
				fail("Zones.Get", err)
			} else {
				ok(fmt.Sprintf("Zones.Get (%s)", got.Name))
			}
			list, _, err := scoped.Zones.List(ctx, &jumpserver.ListOptions{Limit: 15})
			if err != nil {
				fail("Zones.List", err)
			} else {
				ok(fmt.Sprintf("Zones.List (%d)", len(list)))
			}
		}
	}

	// ============================================================
	section("Labels CRUD")
	{
		l, _, err := scoped.Labels.Create(ctx, &model.LabelRequest{Name: "lbl-" + ts, Value: "v1"})
		if err != nil {
			fail("Labels.Create", err)
		} else {
			createdLabelID = l.ID
			ok(fmt.Sprintf("Labels.Create (id=%s)", l.ID))
			u, _, err := scoped.Labels.Update(ctx, l.ID, &model.LabelRequest{Name: "lbl-upd-" + ts, Value: "v2"})
			if err != nil {
				fail("Labels.Update", err)
			} else {
				ok(fmt.Sprintf("Labels.Update (%s)", u.Name))
			}
			got, _, err := scoped.Labels.Get(ctx, l.ID)
			if err != nil {
				fail("Labels.Get", err)
			} else {
				ok(fmt.Sprintf("Labels.Get (%s)", got.Name))
			}
			list, _, err := scoped.Labels.List(ctx, &jumpserver.ListOptions{Limit: 15})
			if err != nil {
				fail("Labels.List", err)
			} else {
				ok(fmt.Sprintf("Labels.List (%d)", len(list)))
			}
		}
	}

	// ============================================================
	section("Hosts CRUD")
	{
		platforms, _, _ := client.Platforms.List(ctx, &jumpserver.ListOptions{Limit: 100})
		var platID int
		for _, p := range platforms {
			if strings.Contains(strings.ToLower(p.Name), "linux") {
				platID = p.ID
				break
			}
		}
		if platID == 0 && len(platforms) > 0 {
			platID = platforms[0].ID
		}
		if platID == 0 {
			skip("Hosts.Create", "no platform")
		} else {
			h, _, err := scoped.Hosts.Create(ctx, &model.AssetRequest{
				Name: "host-" + ts, Address: "192.168.1." + ts[:min(3, len(ts))],
				Platform: platID, Protocols: []model.NamePort{{Name: "ssh", Port: 22}},
			})
			if err != nil {
				fail("Hosts.Create", err)
			} else {
				createdHostID = h.ID
				ok(fmt.Sprintf("Hosts.Create (id=%s)", h.ID))
				u, _, err := scoped.Hosts.Update(ctx, h.ID, &model.AssetRequest{
					Name: "host-upd-" + ts, Address: "10.0.0.1", Platform: platID,
					Protocols: []model.NamePort{{Name: "ssh", Port: 2222}},
				})
				if err != nil {
					fail("Hosts.Update", err)
				} else {
					ok(fmt.Sprintf("Hosts.Update (%s)", u.Name))
				}
				got, _, err := scoped.Hosts.Get(ctx, h.ID)
				if err != nil {
					fail("Hosts.Get", err)
				} else {
					ok(fmt.Sprintf("Hosts.Get (%s)", got.Name))
				}
				list, _, err := scoped.Hosts.List(ctx, nil, &jumpserver.ListOptions{Limit: 15})
				if err != nil {
					fail("Hosts.List", err)
				} else {
					ok(fmt.Sprintf("Hosts.List (%d)", len(list)))
				}
			}
		}
	}

	// ============================================================
	section("Assets (generic)")
	{
		assetList, _, err := client.Assets.List(ctx, nil, &jumpserver.ListOptions{Limit: 15})
		if err != nil {
			fail("Assets.List", err)
		} else {
			ok(fmt.Sprintf("Assets.List (%d)", len(assetList)))
		}
		if len(assetList) > 0 {
			a, _, err := client.Assets.Get(ctx, assetList[0].ID)
			if err != nil {
				fail("Assets.Get", err)
			} else {
				ok(fmt.Sprintf("Assets.Get (%s)", a.Name))
			}
		}
	}

	// ============================================================
	section("Category assets CRUD (Databases, Devices, Webs, Clouds, Customs)")
	{
		// Fetch all platforms to find matching ones per category.
		allPlatforms, _, err := client.Platforms.List(ctx, &jumpserver.ListOptions{Limit: 500})
		if err != nil {
			fail("Platforms.List (for categories)", err)
		} else {
			// Build a lookup: category value -> first matching platform.
			platByCategory := make(map[string]model.Platform)
			for _, p := range allPlatforms {
				cat := p.Category.Value
				if cat != "" {
					if _, exists := platByCategory[cat]; !exists {
						platByCategory[cat] = p
					}
				}
			}

			type catTest struct {
				name    string
				catKey  string
				svc     *assets.CategoryService
				address string
			}
			cats := []catTest{
				{"Databases", "database", scoped.Databases, "10.0.1.1"},
				{"Devices", "device", scoped.Devices, "10.0.2.1"},
				{"Webs", "web", scoped.Webs, "10.0.3.1"},
				{"Clouds", "cloud", scoped.Clouds, "https://cloud.example.com"},
				{"Customs", "custom", scoped.Customs, "10.0.5.1"},
			}

			for _, tc := range cats {
				p, found := platByCategory[tc.catKey]
				if !found {
					skip(tc.name+".Create", fmt.Sprintf("no platform for category %q", tc.catKey))
					continue
				}

				// Build protocols from the platform's protocol list.
				var protocols []model.NamePort
				for _, proto := range p.Protocols {
					if proto.Name != "" {
						protocols = append(protocols, model.NamePort{Name: proto.Name, Port: proto.Port})
					}
				}
				if len(protocols) == 0 {
					protocols = []model.NamePort{{Name: "ssh", Port: 22}}
				}

				// Create
				asset, _, err := tc.svc.Create(ctx, &model.AssetRequest{
					Name:      fmt.Sprintf("%s-%s", strings.ToLower(tc.name), ts),
					Address:   tc.address,
					Platform:  p.ID,
					Protocols: protocols,
				})
				if err != nil {
					fail(tc.name+".Create", err)
					continue
				}
				createdCategoryIDs[tc.name] = asset.ID
				ok(fmt.Sprintf("%s.Create (id=%s, platform=%s/%d)", tc.name, asset.ID, p.Name, p.ID))

				// Update
				u, _, err := tc.svc.Update(ctx, asset.ID, &model.AssetRequest{
					Name:      fmt.Sprintf("%s-upd-%s", strings.ToLower(tc.name), ts),
					Address:   tc.address,
					Platform:  p.ID,
					Protocols: protocols,
				})
				if err != nil {
					fail(tc.name+".Update", err)
				} else {
					ok(fmt.Sprintf("%s.Update (%s)", tc.name, u.Name))
				}

				// Get
				got, _, err := tc.svc.Get(ctx, asset.ID)
				if err != nil {
					fail(tc.name+".Get", err)
				} else {
					ok(fmt.Sprintf("%s.Get (%s)", tc.name, got.Name))
				}

				// List
				list, _, err := tc.svc.List(ctx, nil, &jumpserver.ListOptions{Limit: 15})
				if err != nil {
					fail(tc.name+".List", err)
				} else {
					ok(fmt.Sprintf("%s.List (%d)", tc.name, len(list)))
				}
			}
		}
	}

	// ============================================================
	section("Gateways CRUD")
	{
		platforms, _, _ := client.Platforms.List(ctx, &jumpserver.ListOptions{Limit: 200})
		var gwPlat int
		for _, p := range platforms {
			if strings.Contains(strings.ToLower(p.Name), "gateway") || strings.Contains(strings.ToLower(p.Name), "网关") {
				gwPlat = p.ID
				break
			}
		}
		if gwPlat == 0 {
			skip("Gateways.Create", "no gateway-type platform found")
		} else {
			gw, _, err := scoped.Gateways.Create(ctx, &model.GatewayRequest{
				Name: "gw-" + ts, Address: "172.16.0.1", Platform: gwPlat,
				Protocols: []model.NamePort{{Name: "ssh", Port: 22}},
			})
			if err != nil {
				fail("Gateways.Create", err)
			} else {
				ok(fmt.Sprintf("Gateways.Create (id=%s)", gw.ID))
				u, _, err := scoped.Gateways.Update(ctx, gw.ID, &model.GatewayRequest{
					Name: "gw-upd-" + ts, Address: "172.16.0.2", Platform: gwPlat,
					Protocols: []model.NamePort{{Name: "ssh", Port: 22}},
				})
				if err != nil {
					fail("Gateways.Update", err)
				} else {
					ok(fmt.Sprintf("Gateways.Update (%s)", u.Name))
				}
				got, _, err := scoped.Gateways.Get(ctx, gw.ID)
				if err != nil {
					fail("Gateways.Get", err)
				} else {
					ok(fmt.Sprintf("Gateways.Get (%s)", got.Name))
				}
				_, err = scoped.Gateways.Delete(ctx, gw.ID)
				if err != nil {
					fail("Gateways.Delete", err)
				} else {
					ok("Gateways.Delete")
				}
			}
			list, _, err := scoped.Gateways.List(ctx, &jumpserver.ListOptions{Limit: 15})
			if err != nil {
				fail("Gateways.List", err)
			} else {
				ok(fmt.Sprintf("Gateways.List (%d)", len(list)))
			}
		}
	}

	// ============================================================
	section("AccountTemplates CRUD")
	{
		t, _, err := scoped.AccountTemplates.Create(ctx, &model.AccountTemplateRequest{
			Name: "tpl-" + ts, Username: "testuser", SecretType: "password",
		})
		if err != nil {
			fail("AccountTemplates.Create", err)
		} else {
			createdTemplateID = t.ID
			ok(fmt.Sprintf("AccountTemplates.Create (id=%s)", t.ID))
			u, _, err := scoped.AccountTemplates.Update(ctx, t.ID, &model.AccountTemplateRequest{
				Name: "tpl-upd-" + ts, Username: "testuser2", SecretType: "password",
			})
			if err != nil {
				fail("AccountTemplates.Update", err)
			} else {
				ok(fmt.Sprintf("AccountTemplates.Update (%s)", u.Name))
			}
			got, _, err := scoped.AccountTemplates.Get(ctx, t.ID)
			if err != nil {
				fail("AccountTemplates.Get", err)
			} else {
				ok(fmt.Sprintf("AccountTemplates.Get (%s)", got.Name))
			}
			list, _, err := scoped.AccountTemplates.List(ctx, &jumpserver.ListOptions{Limit: 15})
			if err != nil {
				fail("AccountTemplates.List", err)
			} else {
				ok(fmt.Sprintf("AccountTemplates.List (%d)", len(list)))
			}
		}
	}

	// ============================================================
	section("Accounts CRUD")
	{
		if createdHostID == "" {
			skip("Accounts.Create", "no host")
		} else {
			a, _, err := scoped.Accounts.Create(ctx, &model.AccountRequest{
				Name: "acct-" + ts, Username: "testacct", Asset: createdHostID,
				SecretType: "password", Secret: "TestPass123!",
			})
			if err != nil {
				fail("Accounts.Create", err)
			} else {
				createdAccountID = a.ID
				ok(fmt.Sprintf("Accounts.Create (id=%s)", a.ID))
				u, _, err := scoped.Accounts.Update(ctx, a.ID, &model.AccountRequest{
					Username: "testacct-upd", Asset: createdHostID, SecretType: "password",
				})
				if err != nil {
					fail("Accounts.Update", err)
				} else {
					ok(fmt.Sprintf("Accounts.Update (%s)", u.Username))
				}
				got, _, err := scoped.Accounts.Get(ctx, a.ID)
				if err != nil {
					fail("Accounts.Get", err)
				} else {
					ok(fmt.Sprintf("Accounts.Get (%s)", got.Username))
				}
				list, _, err := scoped.Accounts.List(ctx, &jumpserver.ListOptions{Limit: 15})
				if err != nil {
					fail("Accounts.List", err)
				} else {
					ok(fmt.Sprintf("Accounts.List (%d)", len(list)))
				}
			}
		}
	}

	// ============================================================
	section("ChangeSecrets CRUD")
	{
		cs, _, err := scoped.ChangeSecrets.Create(ctx, &model.ChangeSecretAutomationRequest{
			Name: "cs-" + ts, Accounts: []string{"root"}, SecretType: "password",
			SecretStrategy: "specific", IsPeriodic: true, Interval: 24,
		})
		if err != nil {
			fail("ChangeSecrets.Create", err)
		} else {
			ok(fmt.Sprintf("ChangeSecrets.Create (id=%s)", cs.ID))
			u, _, err := scoped.ChangeSecrets.Update(ctx, cs.ID, &model.ChangeSecretAutomationRequest{
				Name: "cs-upd-" + ts, Accounts: []string{"root"}, SecretType: "password",
				SecretStrategy: "specific", IsPeriodic: true, Interval: 48,
			})
			if err != nil {
				fail("ChangeSecrets.Update", err)
			} else {
				ok(fmt.Sprintf("ChangeSecrets.Update (%s)", u.Name))
			}
			got, _, err := scoped.ChangeSecrets.Get(ctx, cs.ID)
			if err != nil {
				fail("ChangeSecrets.Get", err)
			} else {
				ok(fmt.Sprintf("ChangeSecrets.Get (%s)", got.Name))
			}
			list, _, err := scoped.ChangeSecrets.List(ctx, &jumpserver.ListOptions{Limit: 15})
			if err != nil {
				fail("ChangeSecrets.List", err)
			} else {
				ok(fmt.Sprintf("ChangeSecrets.List (%d)", len(list)))
			}
			_, err = scoped.ChangeSecrets.Delete(ctx, cs.ID)
			if err != nil {
				fail("ChangeSecrets.Delete", err)
			} else {
				ok("ChangeSecrets.Delete")
			}
		}
	}

	// ============================================================
	section("AccountBackups CRUD")
	{
		bp, _, err := scoped.AccountBackups.Create(ctx, &model.AccountBackupPlanRequest{
			Name: "bp-" + ts, Accounts: []string{"root"}, SecretType: "password",
			IsPeriodic: true, Interval: 24,
		})
		if err != nil {
			fail("AccountBackups.Create", err)
		} else {
			ok(fmt.Sprintf("AccountBackups.Create (id=%s)", bp.ID))
			u, _, err := scoped.AccountBackups.Update(ctx, bp.ID, &model.AccountBackupPlanRequest{
				Name: "bp-upd-" + ts, Accounts: []string{"root"}, SecretType: "password",
				IsPeriodic: true, Interval: 48,
			})
			if err != nil {
				fail("AccountBackups.Update", err)
			} else {
				ok(fmt.Sprintf("AccountBackups.Update (%s)", u.Name))
			}
			got, _, err := scoped.AccountBackups.Get(ctx, bp.ID)
			if err != nil {
				fail("AccountBackups.Get", err)
			} else {
				ok(fmt.Sprintf("AccountBackups.Get (%s)", got.Name))
			}
			list, _, err := scoped.AccountBackups.List(ctx, &jumpserver.ListOptions{Limit: 15})
			if err != nil {
				fail("AccountBackups.List", err)
			} else {
				ok(fmt.Sprintf("AccountBackups.List (%d)", len(list)))
			}
			_, err = scoped.AccountBackups.Delete(ctx, bp.ID)
			if err != nil {
				fail("AccountBackups.Delete", err)
			} else {
				ok("AccountBackups.Delete")
			}
		}
	}

	// ============================================================
	section("Permissions CRUD")
	{
		p, _, err := scoped.Permissions.Create(ctx, &model.AssetPermissionRequest{Name: "perm-" + ts})
		if err != nil {
			fail("Permissions.Create", err)
		} else {
			createdPermID = p.ID
			ok(fmt.Sprintf("Permissions.Create (id=%s)", p.ID))
			u, _, err := scoped.Permissions.Update(ctx, p.ID, &model.AssetPermissionRequest{Name: "perm-upd-" + ts})
			if err != nil {
				fail("Permissions.Update", err)
			} else {
				ok(fmt.Sprintf("Permissions.Update (%s)", u.Name))
			}
			got, _, err := scoped.Permissions.Get(ctx, p.ID)
			if err != nil {
				fail("Permissions.Get", err)
			} else {
				ok(fmt.Sprintf("Permissions.Get (%s)", got.Name))
			}
			list, _, err := scoped.Permissions.List(ctx, &jumpserver.ListOptions{Limit: 15})
			if err != nil {
				fail("Permissions.List", err)
			} else {
				ok(fmt.Sprintf("Permissions.List (%d)", len(list)))
			}
			_, err = scoped.Permissions.Delete(ctx, p.ID)
			if err != nil {
				fail("Permissions.Delete", err)
			} else {
				ok("Permissions.Delete")
				createdPermID = ""
			}
		}
	}

	// ============================================================
	section("CommandFilters CRUD")
	{
		cf, _, err := scoped.CommandFilters.Create(ctx, &model.CommandFilterRequest{
			Name: "cf-" + ts, Action: "reject",
			Users: map[string]any{"type": "all"}, Assets: map[string]any{"type": "all"},
			Accounts: []string{"*"},
		})
		if err != nil {
			fail("CommandFilters.Create", err)
		} else {
			createdCmdFilterID = cf.ID
			ok(fmt.Sprintf("CommandFilters.Create (id=%s)", cf.ID))
			u, _, err := scoped.CommandFilters.Update(ctx, cf.ID, &model.CommandFilterRequest{
				Name: "cf-upd-" + ts, Action: "reject",
				Users: map[string]any{"type": "all"}, Assets: map[string]any{"type": "all"},
				Accounts: []string{"*"},
			})
			if err != nil {
				fail("CommandFilters.Update", err)
			} else {
				ok(fmt.Sprintf("CommandFilters.Update (%s)", u.Name))
			}
			got, _, err := scoped.CommandFilters.Get(ctx, cf.ID)
			if err != nil {
				fail("CommandFilters.Get", err)
			} else {
				ok(fmt.Sprintf("CommandFilters.Get (%s)", got.Name))
			}
			list, _, err := scoped.CommandFilters.List(ctx, &jumpserver.ListOptions{Limit: 15})
			if err != nil {
				fail("CommandFilters.List", err)
			} else {
				ok(fmt.Sprintf("CommandFilters.List (%d)", len(list)))
			}
		}
	}

	// ============================================================
	section("CommandGroups CRUD")
	{
		cg, _, err := scoped.CommandFilters.CreateGroup(ctx, &model.CommandGroupRequest{
			Name: "cg-" + ts, Type: map[string]string{"value": "command"}, Content: "rm -rf",
		})
		if err != nil {
			fail("CommandGroups.Create", err)
		} else {
			createdCmdGroupID = cg.ID
			ok(fmt.Sprintf("CommandGroups.Create (id=%s)", cg.ID))
			u, _, err := scoped.CommandFilters.UpdateGroup(ctx, cg.ID, &model.CommandGroupRequest{
				Name: "cg-upd-" + ts, Type: map[string]string{"value": "command"}, Content: "rm -rf /",
			})
			if err != nil {
				fail("CommandGroups.Update", err)
			} else {
				ok(fmt.Sprintf("CommandGroups.Update (%s)", u.Name))
			}
			got, _, err := scoped.CommandFilters.GetGroup(ctx, cg.ID)
			if err != nil {
				fail("CommandGroups.Get", err)
			} else {
				ok(fmt.Sprintf("CommandGroups.Get (%s)", got.Name))
			}
			list, _, err := scoped.CommandFilters.ListGroups(ctx, &jumpserver.ListOptions{Limit: 15})
			if err != nil {
				fail("CommandGroups.List", err)
			} else {
				ok(fmt.Sprintf("CommandGroups.List (%d)", len(list)))
			}
		}
	}

	// ============================================================
	section("LoginACLs")
	{
		list, _, err := client.LoginACLs.List(ctx, &jumpserver.ListOptions{Limit: 15})
		if err != nil {
			fail("LoginACLs.List", err)
		} else {
			ok(fmt.Sprintf("LoginACLs.List (%d)", len(list)))
		}
		if len(list) > 0 {
			a, _, err := client.LoginACLs.Get(ctx, list[0].ID)
			if err != nil {
				fail("LoginACLs.Get", err)
			} else {
				ok(fmt.Sprintf("LoginACLs.Get (%s)", a.Name))
			}
		}
	}

	// ============================================================
	section("Tickets")
	{
		list, _, err := client.Tickets.List(ctx, &jumpserver.ListOptions{Limit: 15})
		if err != nil {
			fail("Tickets.List", err)
		} else {
			ok(fmt.Sprintf("Tickets.List (%d)", len(list)))
		}
		if len(list) > 0 {
			t, _, err := client.Tickets.Get(ctx, list[0].ID)
			if err != nil {
				fail("Tickets.Get", err)
			} else {
				ok(fmt.Sprintf("Tickets.Get (%s)", t.Title))
			}
		}
		flows, _, err := client.Tickets.ListFlows(ctx, &jumpserver.ListOptions{Limit: 15})
		if err != nil {
			fail("Tickets.ListFlows", err)
		} else {
			ok(fmt.Sprintf("Tickets.ListFlows (%d)", len(flows)))
		}
	}

	// ============================================================
	section("Audits")
	{
		sessions, _, err := client.Audits.ListSessions(ctx, &jumpserver.ListOptions{Limit: 3})
		if err != nil {
			fail("Audits.ListSessions", err)
		} else {
			ok(fmt.Sprintf("Audits.ListSessions (%d)", len(sessions)))
		}
		if len(sessions) > 0 {
			s, _, err := client.Audits.GetSession(ctx, sessions[0].ID)
			if err != nil {
				fail("Audits.GetSession", err)
			} else {
				ok(fmt.Sprintf("Audits.GetSession (%s)", s.Asset))
			}
		}
		cmds, _, err := client.Audits.ListCommands(ctx, &jumpserver.ListOptions{Limit: 3})
		if err != nil {
			fail("Audits.ListCommands", err)
		} else {
			ok(fmt.Sprintf("Audits.ListCommands (%d)", len(cmds)))
		}
		ftp, _, err := client.Audits.ListFTPLogs(ctx, &jumpserver.ListOptions{Limit: 3})
		if err != nil {
			fail("Audits.ListFTPLogs", err)
		} else {
			ok(fmt.Sprintf("Audits.ListFTPLogs (%d)", len(ftp)))
		}
		ll, _, err := client.Audits.ListLoginLogs(ctx, &jumpserver.ListOptions{Limit: 3})
		if err != nil {
			fail("Audits.ListLoginLogs", err)
		} else {
			ok(fmt.Sprintf("Audits.ListLoginLogs (%d)", len(ll)))
		}
		ol, _, err := client.Audits.ListOperateLogs(ctx, &jumpserver.ListOptions{Limit: 3})
		if err != nil {
			fail("Audits.ListOperateLogs", err)
		} else {
			ok(fmt.Sprintf("Audits.ListOperateLogs (%d)", len(ol)))
		}
	}

	// ============================================================
	section("Terminal")
	{
		methods, _, err := client.Terminal.ConnectMethods(ctx)
		if err != nil {
			fail("Terminal.ConnectMethods", err)
		} else {
			ok(fmt.Sprintf("Terminal.ConnectMethods (%d keys)", len(methods)))
		}
	}

	// ============================================================
	section("Xpack")
	{
		lic, _, err := client.Xpack.License(ctx)
		if err != nil {
			fail("Xpack.License", err)
		} else {
			ok(fmt.Sprintf("Xpack.License (valid=%v, edition=%s, corporation=%s)", lic.IsValid, lic.Edition, lic.Corporation))
		}
	}

	// ============================================================
	section("WalkPages")
	{
		var all []string
		err := jumpserver.WalkPages(ctx, &jumpserver.ListOptions{Limit: 10}, 20,
			func(ctx context.Context, opts *jumpserver.ListOptions) (*jumpserver.Response, error) {
				usrs, resp, err := client.Users.List(ctx, nil, opts)
				if err != nil {
					return resp, err
				}
				for _, u := range usrs {
					all = append(all, u.Username)
				}
				return resp, nil
			})
		if err != nil {
			fail("WalkPages", err)
		} else {
			ok(fmt.Sprintf("WalkPages (%d users)", len(all)))
		}
	}

	// ============================================================
	section("WithOrgScope")
	{
		s2 := client.WithOrgScope(model.JMSDefaultOrg)
		p, _, err := s2.Users.Profile(ctx)
		if err != nil {
			fail("WithOrgScope(ROOT).Profile", err)
		} else {
			ok(fmt.Sprintf("WithOrgScope(ROOT).Profile (%s)", p.Username))
		}
	}

	// ============================================================
	section("Cleanup")
	{
		for _, c := range []struct {
			name string
			fn   func() (bool, error)
		}{
			{"Account", func() (bool, error) {
				if createdAccountID == "" {
					return true, nil
				}
				_, e := scoped.Accounts.Delete(ctx, createdAccountID)
				return false, e
			}},
			{"AccountTemplate", func() (bool, error) {
				if createdTemplateID == "" {
					return true, nil
				}
				_, e := scoped.AccountTemplates.Delete(ctx, createdTemplateID)
				return false, e
			}},
			{"Host", func() (bool, error) {
				if createdHostID == "" {
					return true, nil
				}
				_, e := scoped.Hosts.Delete(ctx, createdHostID)
				return false, e
			}},
			{"Category:Database", func() (bool, error) {
				id := createdCategoryIDs["Databases"]
				if id == "" {
					return true, nil
				}
				_, e := scoped.Databases.Delete(ctx, id)
				return false, e
			}},
			{"Category:Device", func() (bool, error) {
				id := createdCategoryIDs["Devices"]
				if id == "" {
					return true, nil
				}
				_, e := scoped.Devices.Delete(ctx, id)
				return false, e
			}},
			{"Category:Web", func() (bool, error) {
				id := createdCategoryIDs["Webs"]
				if id == "" {
					return true, nil
				}
				_, e := scoped.Webs.Delete(ctx, id)
				return false, e
			}},
			{"Category:Cloud", func() (bool, error) {
				id := createdCategoryIDs["Clouds"]
				if id == "" {
					return true, nil
				}
				_, e := scoped.Clouds.Delete(ctx, id)
				return false, e
			}},
			{"Category:Custom", func() (bool, error) {
				id := createdCategoryIDs["Customs"]
				if id == "" {
					return true, nil
				}
				_, e := scoped.Customs.Delete(ctx, id)
				return false, e
			}},
			{"CommandGroup", func() (bool, error) {
				if createdCmdGroupID == "" {
					return true, nil
				}
				_, e := scoped.CommandFilters.DeleteGroup(ctx, createdCmdGroupID)
				return false, e
			}},
			{"CommandFilter", func() (bool, error) {
				if createdCmdFilterID == "" {
					return true, nil
				}
				_, e := scoped.CommandFilters.Delete(ctx, createdCmdFilterID)
				return false, e
			}},
			{"Permission", func() (bool, error) {
				if createdPermID == "" {
					return true, nil
				}
				_, e := scoped.Permissions.Delete(ctx, createdPermID)
				return false, e
			}},
			{"Node", func() (bool, error) {
				if createdNodeID == "" {
					return true, nil
				}
				_, e := scoped.Nodes.Delete(ctx, createdNodeID)
				return false, e
			}},
			{"Zone", func() (bool, error) {
				if createdZoneID == "" {
					return true, nil
				}
				_, e := scoped.Zones.Delete(ctx, createdZoneID)
				return false, e
			}},
			{"Label", func() (bool, error) {
				if createdLabelID == "" {
					return true, nil
				}
				_, e := scoped.Labels.Delete(ctx, createdLabelID)
				return false, e
			}},
		} {
			skipped, err := c.fn()
			if err != nil {
				fail("Cleanup."+c.name, err)
			} else if skipped {
				skip("Cleanup."+c.name, "not created")
			} else {
				ok("Cleanup." + c.name)
			}
		}
	}

	// ============================================================
	fmt.Println(strings.Repeat("=", 60))
	fmt.Printf("Results: %d passed, %d failed, %d skipped, %d total\n",
		passed, failed, skipped, passed+failed+skipped)
	if failed > 0 {
		os.Exit(1)
	}
}
