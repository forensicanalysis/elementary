package pluginlib

type Provider interface {
	List() []Plugin
}
