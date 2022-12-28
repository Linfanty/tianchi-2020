package main

import (
    // "fmt"
	"sort"
	// "regexp" 
	// "strconv"
	// mapset "github.com/deckarep/golang-set"
)

// bcj
type unionSet struct {
    rank []int
    set  []int
}

func NewUnionSet(size int) *unionSet {
    buf1 := make([]int, size)
    for i := 0; i < size; i++ {
        buf1[i] = i
    }
    buf2 := make([]int, size)
    for i := 0; i < size; i++ {
        buf2[i] = 1
    }

    return &unionSet{
        rank: buf2,
        set:  buf1,
    }
}

func (set *unionSet) GetSize() int {
    return len(set.set)
}

func (set *unionSet) GetID(p int) (int, error) {
    if p < 0 || p > len(set.set) {
        print("failed to get ID,index is illegal.")
    }

    return set.getRoot(p), nil
}

func (set *unionSet) getRoot(p int) int {
    for p != set.set[p] {
        set.set[p] = set.set[set.set[p]]
        p = set.set[p]
    }
    return p
}

func (set *unionSet) IsConnected(p, q int) (bool, error) {
    if p < 0 || p > len(set.set) || q < 0 || q > len(set.set) {
		print("error: index is illegal.")
    }
    return set.getRoot(set.set[p]) == set.getRoot(set.set[q]), nil
}

func (set *unionSet) Union(p, q int) error {
    if p < 0 || p > len(set.set) || q < 0 || q > len(set.set) {
        print("error: index is illegal.")
    }

    pRoot := set.getRoot(p)
    qRoot := set.getRoot(q)

    if pRoot != qRoot {
        if set.rank[pRoot] < set.rank[qRoot] {
            set.set[pRoot] = qRoot
        } else if set.rank[qRoot] < set.rank[pRoot] {
            set.set[qRoot] = pRoot
        } else { 
            set.set[pRoot] = qRoot
            set.rank[qRoot] += 1
        }
    }
    return nil
}

func (m *manager) Result() map[string][]string {
	return m.res
}

type manager struct {
	// apps         map[string]int            // sidecars of app
	// dependencies map[string]map[string]int // depended services of app and endpoints of each service
	Pilot map[int]int
	Pilot_Data_Name  map[int]string
	res map[string][]string
	Pilot_Phase2 map[int][]string
	// Pilot_Phase2_Set map[int]map[string]bool
	Set_3_ map[int]map[string]bool
	// App_Vec []APP_VEC
	Pilot_Conn  map[int]int
	father map[int]int
	App []APP
}

func NewManager() *manager {
	ret := &manager{
		// apps:         map[string]int{},
		// dependencies: map[string]map[string]int{},
		Pilot:				map[int]int{},
		Pilot_Data_Name:    map[int]string{},
		res:				map[string][]string{},
		Pilot_Phase2:       map[int][]string{},
		// Pilot_Phase2_Set: map[int]map[string]bool{},
		Set_3_:             map[int]map[string]bool{},
		// App_Vec:			[]APP_VEC{},
		Pilot_Conn:         map[int]int{},
		father:				map[int]int{},
		App:				[]APP{},
	}
	return ret
}

func (m *manager) Reset() {
	//m.apps = map[string]int{}
	//m.dependencies = map[string]map[string]int{}
	m.Pilot = map[int]int{}
	m.Pilot_Data_Name = map[int]string{}
	m.res = map[string][]string{}
	m.Pilot_Phase2 = map[int][]string{}
	// m.Pilot_Phase2_Set = map[int]map[string]bool{}
	// m.App_Vec = []APP_VEC{}
	m.Set_3_ = map[int]map[string]bool{}
	m.Pilot_Conn = map[int]int{}
	m.father = map[int]int{}
	m.App = []APP{}
}

type SRV struct {
	SRV_name string
 	SRV_Weight int
}

type APP_VEC struct {
	ROOT int
	Value int
	Srvs []SRV
	Conn int
}

// first srv weight sort
type APP struct {
	APP_Name string
	APP_Name_Int int
	APP_Weight int
	APP_SRV_Weight int
	APP_Srv []SRV
	// APP_Repeat_app map[string]int
}

func (m *manager) initPilots(pilots []string) {
	index := -1

	for _, pilot := range pilots {
		index += 1
		m.Pilot_Data_Name[index] = pilot		
		m.Pilot[index] = 0
		m.Set_3_[index] = make(map[string]bool)
		m.Pilot_Conn[index] = 0

		// print("pilot init ", pilot , "\n")
	}
}

type T_SRV struct {
	Weight int
	Srvs []SRV
	Conn int
}

func (m *manager) fft(App_Vec []APP_VEC, threshhold float64) (App_Vec_2 []APP_VEC) {
	APP_SRV_Set_2 := make(map[int]map[string]bool)
	App_Repeat_2 := make(map[int]map[int]float64)
	var Set_2_ = make(map[int]map[string]bool)

	for i := 0; i < len(App_Vec); i++ {
		APP_SRV_Set_2[ App_Vec[i].ROOT ] = make(map[string]bool)
		App_Repeat_2[ App_Vec[i].ROOT ] = make(map[int]float64)
	}

	for i := 0; i < len(App_Vec); i++ {
		m_srv := App_Vec[i].Srvs
		for j := 0; j < len(m_srv); j++ {
			APP_SRV_Set_2[ App_Vec[i].ROOT ][ m_srv[j].SRV_name ] = true
		}
	}

	for i := 0; i < len(App_Vec); i++ {
		Set_2_[ App_Vec[i].ROOT ] = make(map[string]bool)

		// print("App_Vec[i].ROOT: ", App_Vec[i].ROOT, " App_Vec[i].Value: ", App_Vec[i].Value, "\n") // debug

		for j := i + 1; j < len(App_Vec); j++ {
			total_i := 0
			total_i_repeat := 0

			if len(App_Vec[i].Srvs) < len(App_Vec[j].Srvs) {
				total_i = App_Vec[j].Value
				m_srv := App_Vec[i].Srvs
				
				for k := 0; k < len(m_srv); k++ {
					if APP_SRV_Set_2[ App_Vec[j].ROOT ][ m_srv[k].SRV_name ] {
						total_i_repeat = total_i_repeat + m_srv[k].SRV_Weight
					} else {
						total_i = total_i + m_srv[k].SRV_Weight
					}
				}
			} else {
				total_i = App_Vec[i].Value
				m_srv := App_Vec[j].Srvs
	
				for k := 0; k < len(m_srv); k++ {
					if APP_SRV_Set_2[ App_Vec[i].ROOT ][ m_srv[k].SRV_name ] {
						total_i_repeat = total_i_repeat + m_srv[k].SRV_Weight
					} else {
						total_i = total_i + m_srv[k].SRV_Weight
					}
				}
			}
			
			App_Repeat_2[ App_Vec[i].ROOT ][ App_Vec[j].ROOT ] = float64(total_i_repeat) / float64(total_i)
		}
	}

	bcj_2 := NewUnionSet(5000)
	
	for i := 0; i < len(App_Vec); i++ {
		for j := i + 1; j < len(App_Vec); j++ {
			if App_Repeat_2[ App_Vec[i].ROOT ][ App_Vec[j].ROOT ] > threshhold {
				bcj_2.Union(App_Vec[i].ROOT, App_Vec[j].ROOT)
				// print(App_Vec[i].ROOT, " ", App_Vec[j].ROOT, " ", App_Repeat_2[ App_Vec[i].ROOT ][ App_Vec[j].ROOT ], "\n")
			}
		}
	}

	APP_Map_2 := make(map[int]T_SRV)

	for i := 0; i < len(App_Vec); i++ {
		x := App_Vec[i].ROOT
		y, _ := bcj_2.GetID(x)

		var m_t_srv_2 T_SRV
		m_t_srv_2.Weight = APP_Map_2[y].Weight
		m_t_srv_2.Srvs = APP_Map_2[y].Srvs
		m_t_srv_2.Conn = APP_Map_2[y].Conn + App_Vec[i].Conn
		
		for j := 0; j < len(App_Vec[i].Srvs); j++ {
			if Set_2_[ y ][ App_Vec[i].Srvs[j].SRV_name ] == false {
				m_t_srv_2.Weight += App_Vec[i].Srvs[j].SRV_Weight
				m_t_srv_2.Srvs = append(m_t_srv_2.Srvs , App_Vec[i].Srvs[j])

				Set_2_[ y ][ App_Vec[i].Srvs[j].SRV_name ] = true
			}
		}
		
		APP_Map_2[y] = m_t_srv_2
	}

	for k, v := range APP_Map_2 {
		var m_App_Vec_2 APP_VEC
		m_App_Vec_2.ROOT = k
		m_App_Vec_2.Value = v.Weight
		m_App_Vec_2.Srvs = v.Srvs
		m_App_Vec_2.Conn = v.Conn

		App_Vec_2 = append(App_Vec_2, m_App_Vec_2)
	}

	sort.Slice(App_Vec_2, func(i, j int) bool {
		if App_Vec_2[i].Value > App_Vec_2[j].Value {
			return true
		}
		return false
	})

	print("App_Vec_2 len: ", len(App_Vec_2), "\n")
	// for i := 0; i < len(App_Vec_2); i++ {
	// 	print("App_Vec_2[i].ROOT: ", App_Vec_2[i].ROOT, " App_Vec_2[i].Value: ", App_Vec_2[i].Value, "\n")
	// }

	for i := 0; i < len(m.App); i++ {
		x := m.App[i].APP_Name_Int

		y := m.father[x]
		z, _ := bcj_2.GetID(y)

		m.father[x] = z 
	}
	return App_Vec_2
}

// UpdateAppDependencies merges new apps and dependencies to existing data.
func (m *manager) UpdateAppDependencies_1(apps map[string]int, dependencies map[string]map[string]int) error {
	var tmp_apps = make(map[string]int)

	for app, num := range apps {
		tmp_apps[app] = num
	}

	// sort
	var APP_SRV_Set = make(map[string]map[string]bool)

	App_Repeat := make(map[int]map[int]float64)
	var Set_ = make(map[int]map[string]bool)

	index := -1
	// for app, sidecars := range tmp_apps {
	for app, deps := range dependencies {
		index = index + 1

		var m_app APP
		var srv_weight int
		APP_SRV_Set[app] = make(map[string]bool)

		var m_srv []SRV

		for srv_name, srv_weight_ := range deps{
			var m_srv_part SRV

			m_srv_part.SRV_name = srv_name
			m_srv_part.SRV_Weight = srv_weight_
			m_srv = append(m_srv, m_srv_part)

			APP_SRV_Set[app][srv_name] = true
			srv_weight = srv_weight + srv_weight_
		}
		m_app.APP_Name = app
		m_app.APP_Weight = tmp_apps[app] // sidecars
		m_app.APP_SRV_Weight = srv_weight
		m_app.APP_Srv = m_srv
		// print(m_app.APP_Srv[0].SRV_Weight, "\t")

		// reg := regexp.MustCompile(`[0-9]+`)
		// data := reg.FindString(app)
		// m_app.APP_Name_Int, _ = strconv.Atoi(data)
		m_app.APP_Name_Int = index

		m.App = append(m.App, m_app)
		
		App_Repeat[ index ] = make(map[int]float64)
		Set_[ index ] = make(map[string]bool)
	}

	sort.Slice(m.App, func(i, j int) bool {
		if m.App[i].APP_SRV_Weight > m.App[j].APP_SRV_Weight {
			return true
		}
		return false
	})
	// print("\n\n")
	// for i := 0; i < len(m.App); i++ {
	// 	print(m.App[i].APP_Name, "  ",  m.App[i].APP_SRV_Weight, "\n")
	// }

	// calculate Repeat
	// print("How many apps : ", index, "\n")

	for i := 0; i < len(m.App); i++ {
		for j := i + 1; j < len(m.App); j++ {
			// App_Repeat[App[j].APP_Name_Int] = make(map[int]float64)
			total_i := 0
			total_i_repeat := 0

			if len(m.App[i].APP_Srv) < len(m.App[j].APP_Srv) {
				total_i = m.App[j].APP_SRV_Weight
				m_srv := m.App[i].APP_Srv
	
				for k := 0; k < len(m_srv); k++ {
					if APP_SRV_Set[ m.App[j].APP_Name ][ m_srv[k].SRV_name ] {
						total_i_repeat = total_i_repeat + m_srv[k].SRV_Weight
					} else {
						total_i = total_i + m_srv[k].SRV_Weight
					}
					// print(m_srv[k].SRV_name, "  ", m_srv[k].SRV_Weight, "\n")
					// print(m.App[i].APP_Name, "  ", m.App[i].APP_Weight, "\n")
				}
			} else {
				total_i = m.App[i].APP_SRV_Weight
				m_srv := m.App[j].APP_Srv
	
				for k := 0; k < len(m_srv); k++ {
					if APP_SRV_Set[ m.App[i].APP_Name ][ m_srv[k].SRV_name ] {
						total_i_repeat = total_i_repeat + m_srv[k].SRV_Weight
					} else {
						total_i = total_i + m_srv[k].SRV_Weight
					}
				}
			}
			
			App_Repeat[ m.App[i].APP_Name_Int ][ m.App[j].APP_Name_Int ] = float64(total_i_repeat) / float64(total_i)
			// App_Repeat[ m.App[j].APP_Name_Int ][ m.App[i].APP_Name_Int ] = float64(total_i_repeat) / float64(total_i)

			// if m.App[i].APP_Name_Int == 591 && m.App[j].APP_Name_Int == 961 {
			// 	print(float64(total_i_repeat), "\t", float64(total_i), "\t", m.App[i].APP_SRV_Weight, "\t", App_Repeat[m.App[i].APP_Name_Int][m.App[j].APP_Name_Int], "\t")
			// 	print(len(m.App[i].APP_Srv), "\t",  m.App[i].APP_Name, "\t", len(m.App[j].APP_Srv), "\t",m.App[j].APP_Name, "\n") 
			// }	
		}
	}


	// TODO BFS

	// bcj
	// total_hb := 0
	bcj := NewUnionSet(5000)
	
	for i := 0; i < len(m.App); i++ {
		for j := i + 1; j < len(m.App); j++ {
			if App_Repeat[m.App[i].APP_Name_Int][m.App[j].APP_Name_Int] > 0.98 {
				// total_hb += 1
				bcj.Union(m.App[i].APP_Name_Int, m.App[j].APP_Name_Int)
				// print(App[i].APP_Name_Int, " ", App[j].APP_Name_Int, " ", App_Repeat[App[i].APP_Name_Int][App[j].APP_Name_Int], "\n")
			}
		}
	}
	
	// print(App_Repeat[171][214], "\n")
	// print("total_hb : ", total_hb, "\n")
	// Assign apps to each pilot location
	// Repeat Rate sort

	APP_Map := make(map[int]T_SRV)
	var App_Vec []APP_VEC
	
	for i := 0; i < len(m.App); i++ {
		x := m.App[i].APP_Name_Int
		y, _ := bcj.GetID(x)

		m.father[ x ] = y
		// 2
		var m_t_srv T_SRV

		m_t_srv.Weight = APP_Map[y].Weight
		m_t_srv.Srvs = APP_Map[y].Srvs
		m_t_srv.Conn = APP_Map[y].Conn + m.App[i].APP_Weight
		
		for j := 0; j < len(m.App[i].APP_Srv); j++ {
			if Set_[ y ][ m.App[i].APP_Srv[j].SRV_name ] == false { // not exist
				m_t_srv.Weight += m.App[i].APP_Srv[j].SRV_Weight
				m_t_srv.Srvs = append(m_t_srv.Srvs, m.App[i].APP_Srv[j])

				Set_[ y ][ m.App[i].APP_Srv[j].SRV_name ] = true
			}
		}

		APP_Map[y] = m_t_srv
	}

	for k, v := range APP_Map {
		var m_App_Vec APP_VEC
		m_App_Vec.ROOT = k
		m_App_Vec.Value = v.Weight
		m_App_Vec.Srvs = v.Srvs
		m_App_Vec.Conn = v.Conn

		App_Vec = append(App_Vec, m_App_Vec)
	}

	sort.Slice(App_Vec, func(i, j int) bool {
		if App_Vec[i].Value > App_Vec[j].Value {
			return true
		}
		return false
	})
	
	print("App_Vec[0].ROOT: ", App_Vec[0].Value, " App_Vec[0].Value: ", App_Vec[0].Value, "\n")
	print("App_Vec[1].ROOT: ", App_Vec[1].Value, " App_Vec[1].Value: ", App_Vec[1].Value, "\n")
	print("\n")

	App_Vec = m.fft(App_Vec, 0.90) // -13  -8
	print("App_Vec[0].ROOT: ", App_Vec[0].Value, " App_Vec[0].Value: ", App_Vec[0].Value, "\n")
	print("App_Vec[1].ROOT: ", App_Vec[1].Value, " App_Vec[1].Value: ", App_Vec[1].Value, "\n")
	print("\n")

	App_Vec = m.fft(App_Vec, 0.83) // -13 -7
	print("App_Vec[0].ROOT: ", App_Vec[0].Value, " App_Vec[0].Value: ", App_Vec[0].Value, "\n")
	print("App_Vec[1].ROOT: ", App_Vec[1].Value, " App_Vec[1].Value: ", App_Vec[1].Value, "\n")
	print("\n")

	Root_to_pilot := make(map[int]int)
	type Alpha_List struct {
		K int
		V float64
		W int
	}
	len_m_Pilot := len(m.Pilot)	
	alpha := 0
	bata := 0

	for p_i := 0; p_i < len_m_Pilot; p_i++ {
		Root_to_pilot[ App_Vec[p_i].ROOT ] = p_i
		m.Pilot_Conn[ p_i ] = App_Vec[p_i].Conn

		for j := 0; j < len(App_Vec[ p_i ].Srvs); j++ {
			if m.Set_3_[ p_i ][ App_Vec[ p_i ].Srvs[j].SRV_name ] == false {
				m.Pilot[ p_i ] += App_Vec[ p_i ].Srvs[j].SRV_Weight
				m.Set_3_[ p_i ][ App_Vec[ p_i ].Srvs[j].SRV_name ] = true
			}
		}
		alpha += m.Pilot_Conn[ p_i ]  
		bata += m.Pilot[ p_i ]
	}
	alpha_threshhold := 0.02

	for i := len_m_Pilot; i < len(App_Vec); i++ {
		// alpha_fen := alpha / len_m_Pilot
		// bata_fen := bata / len_m_Pilot
		var alpha_list []Alpha_List

		for p_i := 0; p_i < len_m_Pilot; p_i++ {
			// print(m.Pilot_Conn[p_i], " ", alpha_fen, " ", m.Pilot[p_i], " ", bata_fen, " ", (m.Pilot_Conn[p_i] < alpha_fen && m.Pilot[p_i] < bata_fen), "\n") // debug

			// if m.Pilot_Conn[p_i] < alpha_fen && m.Pilot[p_i] < bata_fen {
				// tmp_W := 0 
				// for j := 0; j < len(App_Vec[i].Srvs); j++ {
				// 	if m.Set_3_[ p_i ][ App_Vec[i].Srvs[j].SRV_name ] {
				// 		tmp_W += App_Vec[i].Srvs[j].SRV_Weight	
				// 	}
				// }
				// print("tmp_W: ", tmp_W, "\n")

				alpha_weight := float64(m.Pilot_Conn[p_i]) * alpha_threshhold + float64(m.Pilot[p_i]) * (1 - alpha_threshhold)
				
				var m_alpha_list Alpha_List
				m_alpha_list.K = p_i
				m_alpha_list.V = alpha_weight
				// m_alpha_list.W = tmp_W 
				alpha_list = append(alpha_list, m_alpha_list) 				
			// }
		}

		sort.Slice(alpha_list, func(i, j int) bool {
			if alpha_list[i].V < alpha_list[j].V {
				return true
			}
			return false
		})
		// for j := 0; j < len(alpha_list); j++ print(j, " ", alpha_list[j].K, " ", alpha_list[j].V, " ", alpha_list[j].W, "\n")

		chosen_index := -1
		if len(alpha_list) >= 1 {
			chosen_index = alpha_list[0].K
		} else {
			Value := int(^uint(0) >> 1)
			for p_i := 0; p_i < len_m_Pilot; p_i++ {
				if m.Pilot[p_i] < Value {
					chosen_index = p_i
					Value = m.Pilot[p_i]
				}
			}
		}
		
		for j := 0; j < len(App_Vec[i].Srvs); j++ {
			if m.Set_3_[ chosen_index ][ App_Vec[i].Srvs[j].SRV_name ] == false {
				m.Pilot[ chosen_index ] += App_Vec[i].Srvs[j].SRV_Weight
				bata += App_Vec[i].Srvs[j].SRV_Weight
				m.Set_3_[ chosen_index ][ App_Vec[i].Srvs[j].SRV_name ] = true
			}
		}

		alpha += App_Vec[i].Conn
		m.Pilot_Conn[chosen_index] += App_Vec[i].Conn
		Root_to_pilot[ App_Vec[i].ROOT ] = chosen_index

		// for p_i := 0; p_i < len(m.Pilot); p_i++ {
		// 	print("Pilot[p_i].Value: ", m.Pilot[p_i], " ", m.Pilot_Conn[p_i],  "\n")
		// }
		// print("\n")
	}

	// result process
	Pilot_Data := make(map[string][]string)

	for i := 0; i < len(m.App); i++ {
		x := m.App[i].APP_Name_Int
		y := m.father[x]

		chosen_index := Root_to_pilot[y]
		Pilot_name := m.Pilot_Data_Name[chosen_index]
		Pilot_Data[ Pilot_name ] = append(Pilot_Data[ Pilot_name ], m.App[i].APP_Name)
	}


	for p_i := 0; p_i < len(m.Pilot); p_i++ {
		print("Pilot[p_i].Value: ", m.Pilot[p_i], " ", m.Pilot_Conn[ p_i ], "\n")
	}
	print("\n")

	m.res = Pilot_Data

	return nil
}

func (m *manager) UpdateAppDependencies_2(apps map[string]int, dependencies map[string]map[string]int) error {
	
	len_m_Pilot := len(m.Pilot)	
	print("How many: ", len(apps), "\n")
	print(" 1.0 / float64(len_m_Pilot): ", 1.0 / float64(len_m_Pilot), "\n")

	Total_Pilot_Srv := 0
	Total_Pilot_App := 0
	for p_i := 0; p_i < len(m.Pilot); p_i++ {
		Total_Pilot_Srv += m.Pilot[ p_i ]
		Total_Pilot_App += m.Pilot_Conn[ p_i ]
	}

	var tmp_apps = make(map[string]int)
	for app, num := range apps {
		tmp_apps[app] = num
	}

	// sort
	APP_SRV_Set := make(map[string]map[string]bool)
	App_Repeat := make(map[int]map[int]float64)
	Set_ := make(map[int]map[string]bool)
	var phase2_App []APP

	index := -1
	for app, deps := range dependencies {
		index = index + 1

		var m_app APP
		var srv_weight int
		APP_SRV_Set[app] = make(map[string]bool)

		var m_srv []SRV

		for srv_name, srv_weight_ := range deps{
			var m_srv_part SRV

			m_srv_part.SRV_name = srv_name
			m_srv_part.SRV_Weight = srv_weight_
			m_srv = append(m_srv, m_srv_part)

			APP_SRV_Set[app][srv_name] = true
			srv_weight = srv_weight + srv_weight_
		}
		m_app.APP_Name = app
		m_app.APP_Weight = tmp_apps[app] // sidecars
		m_app.APP_SRV_Weight = srv_weight
		m_app.APP_Srv = m_srv
		m_app.APP_Name_Int = index

		phase2_App = append(phase2_App, m_app)
		
		App_Repeat[ index ] = make(map[int]float64)
		Set_[ index ] = make(map[string]bool)
	}

	sort.Slice(phase2_App, func(i, j int) bool {
		if phase2_App[i].APP_Weight > phase2_App[j].APP_Weight {
			return true
		}
		return false
	})

	for i := 0; i < len(phase2_App); i++ {
		app := phase2_App[i].APP_Name	
		type Sort_ struct {
			Index int
			Repeat int
		}

		var m_Sort_ []Sort_

		for p_i := 0; p_i < len(m.Pilot); p_i++ {
			var m_m_Sort_ Sort_
			m_m_Sort_.Index = p_i

			for j:= 0; j < len(phase2_App[i].APP_Srv); j++ {
				if m.Set_3_[ p_i ][ phase2_App[i].APP_Srv[j].SRV_name ] {
					m_m_Sort_.Repeat += phase2_App[i].APP_Srv[j].SRV_Weight
				}
			}
			m_Sort_ = append(m_Sort_, m_m_Sort_)
		}

		sort.Slice(m_Sort_, func(i, j int) bool {
			if m_Sort_[i].Repeat > m_Sort_[j].Repeat {
				return true
			}
			return false
		})

		chosen_index := 0
		
		for i := 0; i < len(m_Sort_); i++ {
			print("m_Sort_[i].Index: ", m_Sort_[i].Index, " m.Pilot[ m_Sort_[i].Index] : ", m.Pilot[ m_Sort_[i].Index ],
			" Total_Pilot_Srv: ", Total_Pilot_Srv, 
			" 0.10: ", 1.0 / float64(len_m_Pilot) * float64(Total_Pilot_Srv),  "\n")

			// float64(m.Pilot[ m_Sort_[i].Index ]) > 1.0 / float64(len_m_Pilot) * float64(Total_Pilot_Srv) * 1.1 {

			// print("index: ", m_Sort_[i].Index, " m.Pilot_Conn : ", m.Pilot_Conn[ m_Sort_[i].Index ],
			// " Total_Pilot_App: ", Total_Pilot_App, 
			// " 0.2: ", 0.2 * float64(Total_Pilot_App),  "\n")

			if float64(m.Pilot_Conn[ m_Sort_[i].Index ]) >= 1.0 / float64(len_m_Pilot) * float64(Total_Pilot_App) * 1.1 {
				print("UP hava skip\n")
				continue
			} else {
				chosen_index = m_Sort_[i].Index
				// print("chosen_index: ", chosen_index, "\n")

				m.Pilot_Conn[ chosen_index ] += tmp_apps[app]
				Total_Pilot_App += tmp_apps[app]
				
				for j:= 0; j < len(phase2_App[i].APP_Srv); j++ {
					if m.Set_3_[ chosen_index ][ phase2_App[i].APP_Srv[j].SRV_name ] == false {
						m.Set_3_[ chosen_index ][ phase2_App[i].APP_Srv[j].SRV_name ] = true

						Total_Pilot_Srv += phase2_App[i].APP_Srv[j].SRV_Weight
						m.Pilot[ chosen_index ] += phase2_App[i].APP_Srv[j].SRV_Weight
					} 
				}
				break
			}
		}

		Pilot_name := m.Pilot_Data_Name[chosen_index]
		m.res[ Pilot_name ] = append(m.res[ Pilot_name ], app)
	}
	// for p_i := 0; p_i < len(m.Pilot); p_i++ {
	// 	print("Pilot[p_i].Value: ", m.Pilot[p_i], " m.Pilot_Conn[ p_i ]: ", m.Pilot_Conn[ p_i ], "\n")
	// }
	return nil
}
