package rundeck

import (
	"encoding/xml"
	"errors"
	"reflect"
	"strings"
)

type Job struct {
	XMLName         xml.Name `xml:"job"`
	ID              string   `xml:"id,attr"`
	Name            string   `xml:"name"`
	Group           string   `xml:"group"`
	Project         string   `xml:"project"`
	Description     string   `xml:"description,omitempty"`
	// These two come from Execution output
	AverageDuration int64   `xml:"averageDuration,attr,omitempty"`
	Options         Options `xml:"options,omitempty"`
}

type Options struct {
	XMLName xml.Name
	Options []Option `xml:"option"`
}

type Option struct {
	XMLName xml.Name `xml:"option"`
	Name    string   `xml:"name,attr"`
	Value   string   `xml:"value,attr,omitempty"`
}

type Jobs struct {
	XMLName xml.Name
	Count   int64 `xml:"count,attr"`
	Jobs    []Job `xml:"job"`
}

type RunOptions struct {
	LogLevel  string `qp:"loglevel,omitempty"`
	AsUser    string `qp:"asUser,omitempty"`
	Arguments string `qp:"argString,omitempty"`
}

func (ro *RunOptions) toQueryParams() (u map[string]string) {
	q := make(map[string]string)
	f := reflect.TypeOf(ro).Elem()
	for i := 0; i < f.NumField(); i++ {
		field := f.Field(i)
		tag := field.Tag
		mytag := tag.Get("qp")
		tokens := strings.Split(mytag, ",")
		if len(tokens) == 1 {
			switch tokens[0] {
			case "-":
			//skip
			default:
				k := tokens[0]
				v := reflect.ValueOf(*ro).Field(i).String()
				q[k] = v
			}
		} else {
			switch tokens[1] {
			case "omitempty":
				if tokens[0] == "" {
					// skip
				} else {
					k := tokens[0]
					v := reflect.ValueOf(*ro).Field(i).String()
					q[k] = v
				}
			default:
			//skip
			}
		}
	}
	return q
}

type JobList struct {
	XMLName xml.Name   `xml:"joblist"`
	Job     JobDetails `xml:"job"`
}

type JobDetails struct {
	ID                string          `xml:"id"`
	Name              string          `xml:"name"`
	LogLevel          string          `xml:"loglevel"`
	Description       string          `xml:"description,omitempty"`
	UUID              string          `xml:"uuid"`
	Group             string          `xml:"group"`
	Context           JobContext      `xml:"context"`
	Notification      JobNotification `xml:"notification"`
	MultipleExections bool            `xml:"multipleExecutions"`
	Dispatch          JobDispatch     `xml:"dispatch"`
	NodeFilters       struct {
						  Filter []string `xml:"filter"`
					  } `xml:"nodefilters"`
	Sequence          JobSequence `xml:"sequence"`
}

type JobSequence struct {
	XMLName   xml.Name
	KeepGoing bool           `xml:"keepgoing,attr"`
	Strategy  string         `xml:"strategy,attr"`
	Steps     []SequenceStep `xml:"command"`
}

type SequenceStep struct {
	XMLName        xml.Name
	Description    string      `xml:"description,omitempty"`
	JobRef         *JobRefStep `xml:"jobref,omitempty"`
	NodeStepPlugin *PluginStep `xml:"node-step-plugin,omitempty"`
	StepPlugin     *PluginStep `xml:"step-plugin,omitempty"`
	Exec           *string     `xml:"exec,omitempty"`
	*ScriptStep    `xml:",omitempty"`
}

type ExecStep struct {
	XMLName xml.Name
	string  `xml:"exec,omitempty"`
}

type ScriptStep struct {
	XMLName           xml.Name
	Script            *string `xml:"script,omitempty"`
	ScriptArgs        *string `xml:"scriptargs,omitempty"`
	ScriptFile        *string `xml:"scriptfile,omitempty"`
	ScriptUrl         *string `xml:"scripturl,omitempty"`
	ScriptInterpreter *string `xml:"scriptinterpreter,omitempty"`
}

type PluginStep struct {
	XMLName       xml.Name
	Type          string `xml:"type,attr"`
	Configuration []struct {
		XMLName xml.Name `xml:"entry"`
		Key     string   `xml:"key,attr"`
		Value   string   `xml:"value,attr"`
	} `xml:"configuration>entry,omitempty"`
}

type JobRefStep struct {
	XMLName  xml.Name
	Name     string `xml:"name,attr,omitempty"`
	Group    string `xml:"group,attr,omitempty"`
	NodeStep bool   `xml:"nodeStep,attr,omitempty"`
}

type JobContext struct {
	XMLName xml.Name     `xml:"context"`
	Project string       `xml:"project"`
	Options *[]JobOption `xml:"options>option,omitempty"`
}

type JobOptions struct {
	XMLName xml.Name
	Options []JobOption `xml:"option"`
}

type JobOption struct {
	XMLName      xml.Name `xml:"option"`
	Name         string   `xml:"name,attr"`
	Required     bool     `xml:"required,attr,omitempty"`
	DefaultValue string   `xml:"value,attr,omitempty"`
	Description  string   `xml:"description,omitempty"`
}

type JobNotifications struct {
	Notifications []JobNotification `xml:"notification,omitempty"`
}

type JobNotification struct {
	XMLName   xml.Name   `xml:"notification"`
	OnStart   JobPlugins `xml:"onstart,omitempty"`
	OnSuccess JobPlugins `xml:"onsuccess,omitempty"`
	OnFailure JobPlugins `xml:"onfailure,omitempty"`
}

type JobPlugins struct {
	Plugins []JobPlugin `xml:"plugin,omitempty"`
}

type JobPlugin struct {
	XMLName       xml.Name               `xml:"plugin"`
	PluginType    string                 `xml:"type,attr"`
	Configuration JobPluginConfiguration `xml:"configuration,omitempty"`
}

type JobPluginConfiguration struct {
	XMLName xml.Name                      `xml:"configuration"`
	Entries []JobPluginConfigurationEntry `xml:"entry,omitempty"`
}

type JobPluginConfigurationEntry struct {
	Key   string `xml:"key,attr"`
	Value string `xml:"value,attr,omitempty"`
}

type JobDispatch struct {
	XMLName           xml.Name `xml:"dispatch"`
	ThreadCount       int64    `xml:"threadcount"`
	KeepGoing         bool     `xml:"keepgoing"`
	ExcludePrecedence bool     `xml:"excludePrecendence"`
	RankOrder         string   `xml:"rankOrder"`
}

func (c *RundeckClient) GetJob(id string) (JobList, error) {
	u := make(map[string]string)
	var data JobList
	err := c.Get(&data, "job/"+id, u)
	return data, err
}

func (c *RundeckClient) GetRequiredOpts(j string) (map[string]string, error) {
	u := make(map[string]string)
	var data JobList
	err := c.Get(&data, "job/"+j, u)
	if err != nil {
		return u, err
	} else {
		if data.Job.Context.Options != nil {
			for _, o := range *data.Job.Context.Options {
				if o.Required {
					if o.DefaultValue == "" {
						u[o.Name] = "<no default>"
					} else {
						u[o.Name] = o.DefaultValue
					}
				}
			}
		}
		return u, nil
	}
}
func (c *RundeckClient) GetOpts(j string) (map[string]string, error) {
	u := make(map[string]string)
	var data JobList
	err := c.Get(&data, "job/"+j, u)
	if err != nil {
		return u, err
	} else {
		if data.Job.Context.Options != nil {
			for _, option := range *data.Job.Context.Options {
				var optionRequirement string
				if option.Required {
					optionRequirement = "* "
				}else {
					optionRequirement = "  "
				}

				var optionValue string
				if option.DefaultValue == "" {
					optionValue =  "<no default>"
				} else {
					optionValue = option.DefaultValue
				}

				u[optionRequirement + option.Name] =  optionValue
			}
		}
		return u, nil
	}
}
func (c *RundeckClient) RunJob(id string, options RunOptions) (Executions, error) {
	u := options.toQueryParams()
	var data Executions

	err := c.Get(&data, "job/"+id+"/run", u)
	return data, err
}

func (c *RundeckClient) FindJobByName(name string, project string) (*JobDetails, error) {
	var job *JobDetails
	var err error
	jobs, err := c.ListJobs(project)
	if err != nil {
		//
	} else {
		if len(jobs.Jobs) > 0 {
			for _, d := range jobs.Jobs {
				if d.Name == name {
					joblist, err := c.GetJob(d.ID)
					if err != nil {
						//
					} else {
						job = &joblist.Job
					}
				}
			}
		} else {
			err = errors.New("No matches found")
		}
	}
	return job, err
}

func (c *RundeckClient) ListJobs(projectId string) (Jobs, error) {
	options := make(map[string]string)
	options["project"] = projectId
	var data Jobs
	err := c.Get(&data, "jobs", options)
	return data, err
}
